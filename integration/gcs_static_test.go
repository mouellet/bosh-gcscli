/*
 * Copyright 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package integration

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/bosh-gcscli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func() {
	Context("static credentials configuration", func() {
		regional := os.Getenv(RegionalBucketEnv)
		multiRegional := os.Getenv(MultiRegionalBucketEnv)

		var ctx AssertContext
		BeforeEach(func() {
			Expect(regional).ToNot(BeEmpty(),
				fmt.Sprintf(NoBucketMsg, RegionalBucketEnv))
			Expect(multiRegional).ToNot(BeEmpty(),
				fmt.Sprintf(NoBucketMsg, MultiRegionalBucketEnv))

			ctx = NewAssertContext(AsStaticCredentials)
		})
		AfterEach(func() {
			ctx.Cleanup()
		})

		configurations := []TableEntry{
			Entry("MultiRegional bucket, default StorageClass", &config.GCSCli{
				BucketName: multiRegional,
			}),
			Entry("Regional bucket, default StorageClass", &config.GCSCli{
				BucketName: regional,
			}),
			Entry("MultiRegional bucket, explicit StorageClass", &config.GCSCli{
				BucketName:   multiRegional,
				StorageClass: "MULTI_REGIONAL",
			}),
			Entry("Regional bucket, explicit StorageClass", &config.GCSCli{
				BucketName:   regional,
				StorageClass: "REGIONAL",
			}),
		}

		encryptedConfigs := []TableEntry{
			Entry("MultiRegional bucket, default StorageClass, encrypted", &config.GCSCli{
				BucketName:    multiRegional,
				EncryptionKey: encryptionKeyBytes,
			}),
			Entry("Regional bucket, default StorageClass, encrypted", &config.GCSCli{
				BucketName:    regional,
				EncryptionKey: encryptionKeyBytes,
			}),
		}
		configurations = append(configurations, encryptedConfigs...)

		DescribeTable("Blobstore lifecycle works",
			func(config *config.GCSCli) {
				ctx.AddConfig(config)
				AssertLifecycleWorks(gcsCLIPath, ctx)
			},
			configurations...)

		DescribeTable("Invalid Delete works",
			func(config *config.GCSCli) {
				ctx.AddConfig(config)
				AssertDeleteNonexistentWorks(gcsCLIPath, ctx)
			},
			configurations...)

		DescribeTable("Multipart Put works",
			func(config *config.GCSCli) {
				ctx.AddConfig(config)
				AssertMultipartPutWorks(gcsCLIPath, ctx)
			},
			configurations...)

		DescribeTable("Invalid Put should fail",
			func(config *config.GCSCli) {
				ctx.AddConfig(config)
				AssertBrokenSourcePutFails(gcsCLIPath, ctx)
			},
			configurations...)

		DescribeTable("Invalid Get should fail",
			func(config *config.GCSCli) {
				ctx.AddConfig(config)
				AssertGetNonexistentFails(gcsCLIPath, ctx)
			},
			configurations...)
	})
})