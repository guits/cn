/*
 * Ceph Nano (C) 2018 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * Below main package has canonical imports for 'go get' and 'go build'
 * to work with all other clones of github.com/ceph/cn repository. For
 * more information refer https://golang.org/doc/go1.4#canonicalimports
 */

package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

// cliS3CmdSync is the Cobra CLI call
func cliS3CmdSync() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync [CLUSTER] [LOCAL_DIR] [BUCKET]",
		Short: "Synchronize a directory tree to S3",
		Args:  cobra.ExactArgs(3),
		Run:   S3CmdSync,
	}
	cmd.Flags().BoolVarP(&debugS3, "debug", "d", false, "Run S3 commands in debug mode")

	return cmd
}

// S3CmdSync wraps s3cmd command in the container
func S3CmdSync(cmd *cobra.Command, args []string) {
	containerNameToShow := args[0]
	containerName := containerNamePrefix + containerNameToShow

	notExistCheck(containerName)
	notRunningCheck(containerName)
	localDir := args[1]
	bucketName := args[2]
	dir := dockerInspect(containerName, "Binds")
	destDir := tempPath

	if localDir != dir {
		destDir = dir + "/" + localDir
		err := copyDir(localDir, destDir)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("Syncing directory '%s' in the '%s' bucket. \n"+
		"It might take some time depending on the amount of data. \n"+
		"Do not expect any output until the upload is finished. \n \n", localDir, bucketName)

	command := []string{"s3cmd", "sync", destDir, "s3://" + bucketName}
	if debugS3 {
		command = append(command, "--debug")
	}

	output := strings.TrimSuffix(string(execContainer(containerName, command)), "\n") + " on cluster " + containerNameToShow
	fmt.Println(output)
}
