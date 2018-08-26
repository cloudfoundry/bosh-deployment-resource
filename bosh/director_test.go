package bosh_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-deployment-resource/bosh"
	"github.com/cloudfoundry/bosh-deployment-resource/bosh/boshfakes"
	"github.com/cloudfoundry/bosh-deployment-resource/concourse"
	"github.com/cppforlife/go-patch/patch"
	"github.com/cppforlife/go-semi-semantic/version"

	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshdirfakes "github.com/cloudfoundry/bosh-cli/director/directorfakes"
	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
)

var _ = Describe("BoshDirector", func() {
	var (
		director         bosh.BoshDirector
		commandRunner    *boshfakes.FakeRunner
		sillyBytes       = []byte{0xFE, 0xED, 0xDE, 0xAD, 0xBE, 0xEF}
		fakeBoshDirector *boshdirfakes.FakeDirector
		loggerOutput     bytes.Buffer
	)

	BeforeEach(func() {
		commandRunner = new(boshfakes.FakeRunner)
		fakeBoshDirector = new(boshdirfakes.FakeDirector)

		director = bosh.NewBoshDirector(
			concourse.Source{Deployment: "cool-deployment"},
			commandRunner,
			fakeBoshDirector,
			&loggerOutput,
		)
	})

	Describe("Deploy", func() {
		It("tells BOSH to deploy the given manifest and parameters", func() {
			vars := map[string]interface{}{"foo": "bar"}
			varKVs := []boshtpl.VarKV{
				{
					Name:  "foo",
					Value: "bar",
				},
			}
			varsFileContents := properYaml(`
				baz: "best-bar"
			`)
			varsFile, _ := ioutil.TempFile("", "var-file-1")
			varsFile.Write(varsFileContents)

			varFile, _ := ioutil.TempFile("", "var-file-key2")
			varFile.Write([]byte("val2"))

			opsFileContents := properYaml(`
				- type: replace
				  path: /my?/new_key
				  value: awesome
			`)
			opsFile, _ := ioutil.TempFile("", "ops-file-1")
			opsFile.Write(opsFileContents)

			noRedact := true
			dryRun := false
			err := director.Deploy(sillyBytes, bosh.DeployParams{
				NoRedact:  noRedact,
				DryRun:    dryRun,
				Vars:      vars,
				VarFiles:  map[string]string{"key2": varFile.Name()},
				VarsFiles: []string{varsFile.Name()},
				OpsFiles:  []string{opsFile.Name()},
			})
			Expect(err).ToNot(HaveOccurred())

			Expect(commandRunner.ExecuteCallCount()).To(Equal(1))

			deployOpts := commandRunner.ExecuteArgsForCall(0).(*boshcmd.DeployOpts)
			Expect(deployOpts.Args.Manifest.Bytes).To(Equal(sillyBytes))
			Expect(deployOpts.NoRedact).To(Equal(noRedact))
			Expect(deployOpts.DryRun).To(Equal(dryRun))
			Expect(deployOpts.VarKVs).To(Equal(varKVs))
			Expect(len(deployOpts.VarsFiles)).To(Equal(1))
			Expect(deployOpts.VarsFiles[0].Vars).To(Equal(boshtpl.StaticVariables{
				"baz": "best-bar",
			}))

			Expect(len(deployOpts.VarFiles)).To(Equal(1))
			Expect(deployOpts.VarFiles[0].Vars).To(Equal(boshtpl.StaticVariables{
				"key2": "val2",
			}))
			Expect(len(deployOpts.OpsFiles)).To(Equal(1))

			pathPointer, _ := patch.NewPointerFromString("/my?/new_key")
			Expect(deployOpts.OpsFiles[0].Ops).To(Equal(patch.Ops{
				patch.ReplaceOp{
					Path:  pathPointer,
					Value: "awesome",
				},
			}))
		})

		Describe("VarsStore", func() {
			Context("when one is provided", func() {
				It("is used", func() {
					varsStore, _ := ioutil.TempFile("", "vars-store")
					err := director.Deploy(sillyBytes, bosh.DeployParams{
						VarsStore: varsStore.Name(),
					})
					Expect(err).ToNot(HaveOccurred())

					deployOpts := commandRunner.ExecuteArgsForCall(0).(*boshcmd.DeployOpts)
					Expect(deployOpts.VarsFSStore.FS).ToNot(BeNil())
				})
			})

			Context("when one is not provided", func() {
				It("is not used", func() {
					err := director.Deploy(sillyBytes, bosh.DeployParams{})
					Expect(err).ToNot(HaveOccurred())

					deployOpts := commandRunner.ExecuteArgsForCall(0).(*boshcmd.DeployOpts)
					Expect(deployOpts.VarsFSStore.FS).To(BeNil())
				})
			})
		})

		Context("when deploying fails", func() {
			It("returns an error", func() {
				commandRunner.ExecuteReturns(errors.New("Your deploy failed"))

				err := director.Deploy(sillyBytes, bosh.DeployParams{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Your deploy failed"))
			})
		})

		Context("when cleanup is specified", func() {
			It("runs a cleanup after the deploy", func() {
				err := director.Deploy(sillyBytes, bosh.DeployParams{Cleanup: true})
				Expect(err).ToNot(HaveOccurred())

				Expect(commandRunner.ExecuteCallCount()).To(Equal(2))

				deployOpts := commandRunner.ExecuteArgsForCall(0).(*boshcmd.DeployOpts)
				Expect(deployOpts.Args.Manifest.Bytes).To(Equal(sillyBytes))

				cleanUpOpts := commandRunner.ExecuteArgsForCall(1).(*boshcmd.CleanUpOpts)
				Expect(cleanUpOpts.All).To(Equal(false))
			})
		})

		Context("when dryrun is specified", func() {
			It("use dry-run flags", func() {
				dryRun := true
				err := director.Deploy(sillyBytes, bosh.DeployParams{
					DryRun: dryRun,
				})
				Expect(err).ToNot(HaveOccurred())

				Expect(commandRunner.ExecuteCallCount()).To(Equal(1))

				deployOpts := commandRunner.ExecuteArgsForCall(0).(*boshcmd.DeployOpts)
				Expect(deployOpts.Args.Manifest.Bytes).To(Equal(sillyBytes))
				Expect(deployOpts.DryRun).To(Equal(dryRun))
			})
		})

		Context("when skipdrain is specified", func() {
			It("uses skip-drain flag", func() {
				err := director.Deploy(sillyBytes, bosh.DeployParams{
					SkipDrain: []string{"*"},
				})
				Expect(err).ToNot(HaveOccurred())

				Expect(commandRunner.ExecuteCallCount()).To(Equal(1))

				skipDrain := []boshdir.SkipDrain{
					{All: true},
				}

				deployOpts := commandRunner.ExecuteArgsForCall(0).(*boshcmd.DeployOpts)
				Expect(deployOpts.Args.Manifest.Bytes).To(Equal(sillyBytes))
				Expect(deployOpts.SkipDrain).To(Equal(skipDrain))
			})
		})
	})
	Describe("Delete", func() {
		It("tells BOSH to delete the configured deployment", func() {
			err := director.Delete(true)
			Expect(err).ToNot(HaveOccurred())

			Expect(commandRunner.ExecuteCallCount()).To(Equal(1))

			deleteOpts := commandRunner.ExecuteArgsForCall(0).(*boshcmd.DeleteDeploymentOpts)
			Expect(deleteOpts).To(Equal(&boshcmd.DeleteDeploymentOpts{Force: true}))
		})

		Context("when delete fails", func() {
			BeforeEach(func() {
				commandRunner.ExecuteReturns(errors.New("Delete failed!"))
			})

			It("returns the error", func() {
				err := director.Delete(true)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("Delete failed!"))
			})
		})
	})

	Describe("Interpolate", func() {
		var interpolatedBytes = []byte{0xFE, 0xED, 0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xED, 0xDE, 0xAD, 0xBE, 0xEF}

		BeforeEach(func() {
			commandRunner.ExecuteWithWriterStub = func(commandOpts interface{}, writer io.Writer) error {
				_, err := writer.Write(interpolatedBytes)
				Expect(err).NotTo(HaveOccurred())
				return nil
			}
		})

		It("tells interpolates a BOSH manifest", func() {
			vars := map[string]interface{}{"foo": "bar"}
			varKVs := []boshtpl.VarKV{
				{
					Name:  "foo",
					Value: "bar",
				},
			}
			varFileContents := properYaml(`
				baz: "best-bar"
			`)
			varFile, _ := ioutil.TempFile("", "var-file-1")
			varFile.Write(varFileContents)

			opsFileContents := properYaml(`
				- type: replace
				  path: /my?/new_key
				  value: awesome
			`)
			opsFile, _ := ioutil.TempFile("", "ops-file-1")
			opsFile.Write(opsFileContents)

			manifest, err := director.Interpolate(sillyBytes, bosh.InterpolateParams{
				Vars:      vars,
				VarsFiles: []string{varFile.Name()},
				OpsFiles:  []string{opsFile.Name()},
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(manifest).To(Equal(interpolatedBytes))

			Expect(commandRunner.ExecuteWithWriterCallCount()).To(Equal(1))

			opts, _ := commandRunner.ExecuteWithWriterArgsForCall(0)
			interpolateOpts := opts.(*boshcmd.InterpolateOpts)
			Expect(interpolateOpts.Args.Manifest.Bytes).To(Equal(sillyBytes))
			Expect(interpolateOpts.VarKVs).To(Equal(varKVs))
			Expect(len(interpolateOpts.VarsFiles)).To(Equal(1))
			Expect(interpolateOpts.VarsFiles[0].Vars).To(Equal(boshtpl.StaticVariables{
				"baz": "best-bar",
			}))
			Expect(len(interpolateOpts.OpsFiles)).To(Equal(1))

			pathPointer, _ := patch.NewPointerFromString("/my?/new_key")
			Expect(interpolateOpts.OpsFiles[0].Ops).To(Equal(patch.Ops{
				patch.ReplaceOp{
					Path:  pathPointer,
					Value: "awesome",
				},
			}))
		})

		Context("when interpolating fails", func() {
			It("returns an error", func() {
				commandRunner.ExecuteWithWriterReturns(errors.New("Your interpolate failed"))

				_, err := director.Interpolate(sillyBytes, bosh.InterpolateParams{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Your interpolate failed"))
			})
		})
	})

	Describe("DownloadManifest", func() {
		It("gets the deployment manifest", func() {
			fakeDeployment := boshdirfakes.FakeDeployment{}
			fakeBoshDirector.FindDeploymentReturns(&fakeDeployment, nil)
			fakeDeployment.ManifestReturns(string(sillyBytes), nil)

			manifestBytes, err := director.DownloadManifest()
			Expect(err).ToNot(HaveOccurred())

			Expect(manifestBytes).To(Equal(sillyBytes))
			Expect(fakeBoshDirector.FindDeploymentArgsForCall(0)).To(Equal("cool-deployment"))
		})

		Context("when getting the deployment fails", func() {
			It("returns an error", func() {
				fakeDeployment := boshdirfakes.FakeDeployment{}
				fakeBoshDirector.FindDeploymentReturns(&fakeDeployment, errors.New("Your deployment is missing"))

				_, err := director.DownloadManifest()
				Expect(err).To(MatchError(ContainSubstring("Your deployment is missing")))
			})
		})

		Context("when getting the manifest fails", func() {
			It("returns an error", func() {
				fakeDeployment := boshdirfakes.FakeDeployment{}
				fakeBoshDirector.FindDeploymentReturns(&fakeDeployment, nil)
				fakeDeployment.ManifestReturns(string(sillyBytes), errors.New("Your manifest could not be found"))

				_, err := director.DownloadManifest()
				Expect(err).To(MatchError(ContainSubstring("Your manifest could not be found")))
			})
		})
	})

	Describe("UploadRelease", func() {
		It("uploads the given release", func() {
			err := director.UploadRelease("my-cool-release")
			Expect(err).ToNot(HaveOccurred())

			Expect(commandRunner.ExecuteCallCount()).To(Equal(1))

			uploadReleaseOpts := commandRunner.ExecuteArgsForCall(0).(*boshcmd.UploadReleaseOpts)
			Expect(string(uploadReleaseOpts.Args.URL)).To(Equal("my-cool-release"))
		})

		Context("when uploading the release fails", func() {
			It("returns an error", func() {
				commandRunner.ExecuteReturns(errors.New("failed communicating with director"))

				err := director.UploadRelease("my-cool-release")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Could not upload release my-cool-release: failed communicating with director"))
			})
		})
	})

	Describe("ExportReleases", func() {
		fakeDeployment := new(boshdirfakes.FakeDeployment)
		var fakeDeploymentStemcell *boshdirfakes.FakeStemcell
		var fakeDirectorStemcell *boshdirfakes.FakeStemcell

		BeforeEach(func() {
			version1, err := version.NewVersionFromString("123.45")
			Expect(err).ToNot(HaveOccurred())
			version2, err := version.NewVersionFromString("987.65")
			Expect(err).ToNot(HaveOccurred())
			version3, err := version.NewVersionFromString("abc.de")
			Expect(err).ToNot(HaveOccurred())
			stemcellVersion, err := version.NewVersionFromString("3.4.0")
			Expect(err).ToNot(HaveOccurred())

			fakeRelease1 := new(boshdirfakes.FakeRelease)
			fakeRelease1.NameReturns("cool-release")
			fakeRelease1.VersionReturns(version1)

			fakeRelease2 := new(boshdirfakes.FakeRelease)
			fakeRelease2.NameReturns("awesome-release")
			fakeRelease2.VersionReturns(version2)

			fakeRelease3 := new(boshdirfakes.FakeRelease)
			fakeRelease3.NameReturns("not-requested")
			fakeRelease3.VersionReturns(version3)

			fakeDeployment.ReleasesReturns([]boshdir.Release{fakeRelease1, fakeRelease2, fakeRelease3}, nil)

			fakeDeploymentStemcell = new(boshdirfakes.FakeStemcell)
			fakeDeploymentStemcell.NameReturns("bosh-monkey-minix-go_agent")
			fakeDeploymentStemcell.VersionReturns(stemcellVersion)
			fakeDeployment.StemcellsReturns([]boshdir.Stemcell{fakeDeploymentStemcell}, nil)

			fakeDirectorStemcell = new(boshdirfakes.FakeStemcell)
			fakeDirectorStemcell.NameReturns("bosh-monkey-minix-go_agent")
			fakeDirectorStemcell.OSNameReturns("minix")
			fakeDirectorStemcell.VersionReturns(stemcellVersion)
			fakeBoshDirector.StemcellsReturns([]boshdir.Stemcell{fakeDirectorStemcell}, nil)

			fakeBoshDirector.FindDeploymentReturns(fakeDeployment, nil)
		})

		It("downloads the given releases", func() {
			err := director.ExportReleases("/tmp/foo", []bosh.ReleaseSpec{
				{Name: "cool-release"},
				{
					Name: "awesome-release",
					Jobs: []string{
						"nice-job",
						"well-done",
					},
				},
			})
			Expect(err).ToNot(HaveOccurred())

			Expect(fakeBoshDirector.FindDeploymentCallCount()).To(Equal(1))
			Expect(fakeBoshDirector.FindDeploymentArgsForCall(0)).To(Equal("cool-deployment"))

			Expect(commandRunner.ExecuteWithDefaultOverrideCallCount()).To(Equal(2))

			opts, optFunc, _ := commandRunner.ExecuteWithDefaultOverrideArgsForCall(0)
			exportReleaseOpts, _ := opts.(*boshcmd.ExportReleaseOpts)
			Expect(string(exportReleaseOpts.Args.ReleaseSlug.Name())).To(Equal("cool-release"))
			Expect(string(exportReleaseOpts.Args.ReleaseSlug.Version())).To(Equal("123.45"))
			Expect(string(exportReleaseOpts.Args.OSVersionSlug.OS())).To(Equal("minix"))
			Expect(string(exportReleaseOpts.Args.OSVersionSlug.Version())).To(Equal("3.4.0"))

			fixedOpts, err := optFunc(&boshcmd.ExportReleaseOpts{Directory: boshcmd.DirOrCWDArg{Path: "wrong-path"}})
			Expect(err).ToNot(HaveOccurred())
			Expect(fixedOpts.(*boshcmd.ExportReleaseOpts).Directory.Path).To(Equal("/tmp/foo"))

			opts, optFunc, _ = commandRunner.ExecuteWithDefaultOverrideArgsForCall(1)
			exportReleaseOpts, _ = opts.(*boshcmd.ExportReleaseOpts)
			Expect(string(exportReleaseOpts.Args.ReleaseSlug.Name())).To(Equal("awesome-release"))
			Expect(string(exportReleaseOpts.Args.ReleaseSlug.Version())).To(Equal("987.65"))
			Expect(string(exportReleaseOpts.Args.OSVersionSlug.OS())).To(Equal("minix"))
			Expect(string(exportReleaseOpts.Args.OSVersionSlug.Version())).To(Equal("3.4.0"))
			Expect(exportReleaseOpts.Jobs).To(Equal([]string{
				"nice-job",
				"well-done",
			}))
		})

		Context("when requesting a release not in the manifest", func() {
			It("errors before downloading any releases", func() {
				err := director.ExportReleases("/tmp/foo", []bosh.ReleaseSpec{
					{Name: "cool-release"},
					{Name: "awesome-release"},
					{Name: "missing-release"},
				})
				Expect(err).To(MatchError(ContainSubstring("could not find release missing-release")))

				Expect(commandRunner.ExecuteCallCount()).To(Equal(0))
			})
		})

		Context("when there is more than one stemcell in the manifest", func() {
			It("errors before downloading any releases", func() {
				fakeDeployment.StemcellsReturns([]boshdir.Stemcell{fakeDeploymentStemcell, fakeDeploymentStemcell}, nil)
				err := director.ExportReleases("/tmp/foo", []bosh.ReleaseSpec{
					{Name: "cool-release"},
					{Name: "awesome-release"},
				})
				Expect(err).To(MatchError(ContainSubstring("exporting releases from a deployment with multiple stemcells is unsupported")))

				Expect(commandRunner.ExecuteCallCount()).To(Equal(0))
			})
		})

		Context("when getting the deployment fails", func() {
			It("returns an error", func() {
				fakeBoshDirector.FindDeploymentReturns(fakeDeployment, errors.New("foo"))

				err := director.ExportReleases("/tmp/foo", []bosh.ReleaseSpec{
					{Name: "cool-release"},
				})
				Expect(err).To(MatchError(ContainSubstring("could not export releases: could not fetch deployment cool-deployment: foo")))
			})
		})

		Context("when exporting releases fails", func() {
			It("returns an error", func() {
				commandRunner.ExecuteWithDefaultOverrideReturns(errors.New("failed communicating with director"))

				err := director.ExportReleases("/tmp/foo", []bosh.ReleaseSpec{
					{Name: "cool-release"},
				})
				Expect(err).To(MatchError(ContainSubstring("could not export release cool-release: failed communicating with director")))
			})
		})

		Context("when getting releases fails", func() {
			It("returns an error", func() {
				fakeDeployment.ReleasesReturns([]boshdir.Release{}, errors.New("foo"))

				err := director.ExportReleases("/tmp/foo", []bosh.ReleaseSpec{
					{Name: "cool-release"},
				})
				Expect(err).To(MatchError(ContainSubstring("could not export releases: could not fetch releases: foo")))
			})
		})

		Context("when getting stemcells fails", func() {
			Context("from the deployment", func() {
				It("returns an error", func() {
					fakeDeployment.StemcellsReturns([]boshdir.Stemcell{}, errors.New("foo"))

					err := director.ExportReleases("/tmp/foo", []bosh.ReleaseSpec{
						{Name: "cool-release"},
					})
					Expect(err).To(MatchError(ContainSubstring("could not export releases: could not fetch stemcells: foo")))
				})
			})

			Context("from the director", func() {
				It("returns an error", func() {
					fakeBoshDirector.StemcellsReturns([]boshdir.Stemcell{}, errors.New("foo"))

					err := director.ExportReleases("/tmp/foo", []bosh.ReleaseSpec{
						{Name: "cool-release"},
					})
					Expect(err).To(MatchError(ContainSubstring("could not export releases: could not fetch stemcells: foo")))
				})
			})

		})
	})

	Describe("UploadStemcell", func() {
		It("uploads the given stemcell", func() {
			err := director.UploadStemcell("my-cool-stemcell")
			Expect(err).ToNot(HaveOccurred())

			Expect(commandRunner.ExecuteCallCount()).To(Equal(1))

			uploadStemcellOpts := commandRunner.ExecuteArgsForCall(0).(*boshcmd.UploadStemcellOpts)
			Expect(string(uploadStemcellOpts.Args.URL)).To(Equal("my-cool-stemcell"))
		})

		Context("when uploading the stemcell fails", func() {
			It("returns an error", func() {
				commandRunner.ExecuteReturns(errors.New("failed communicating with director"))

				err := director.UploadStemcell("my-cool-stemcell")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Could not upload stemcell my-cool-stemcell: failed communicating with director"))
			})
		})
	})

	Describe("WaitForDeployLock", func() {
		It("waits for the lock to be released", func() {
			err := director.WaitForDeployLock()
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeBoshDirector.LocksCallCount()).To(Equal(1))
		})

		Context("when there are locks", func() {
			BeforeEach(func() {
				fakeBoshDirector.LocksReturns(
					[]boshdir.Lock{
						{Resource: []string{"other-identifier", "not-my-deployment"}},
						{Resource: []string{"other-identifier", "cool-deployment"}},
					},
					nil,
				)

				fakeBoshDirector.LocksReturnsOnCall(
					1,
					[]boshdir.Lock{
						{Resource: []string{"other-identifier", "not-my-deployment"}},
					},
					nil,
				)
			})

			It("waits for the lock to be released", func() {
				err := director.WaitForDeployLock()
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeBoshDirector.LocksCallCount()).To(Equal(2))

				//logs output so the user knows what is happening
				Expect(loggerOutput.String()).To(ContainSubstring("Waiting for deployment lock. Done\n"))
			})
		})

		Context("when checking the lock fails", func() {
			BeforeEach(func() {
				fakeBoshDirector.LocksReturns([]boshdir.Lock{{Resource: []string{}}}, errors.New("Failed to fetch locks"))
			})

			It("returns an error", func() {
				err := director.WaitForDeployLock()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
