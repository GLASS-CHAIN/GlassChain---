# dapp ci test cases compose rule/ ci test cases

## Framework introduction
Note: The dapp keyword here refers to the name of each plug-in, such as relay, paracross
The testcase of each dapp and the files related to CI testing, such as Dockerfile, are deployed in the cmd/build directory under their own directory.
For example, paracross is under plugin/dapp/paracross/cmd/build

There are two files in the cmd directory, Makefile and build.sh, which are responsible for copying the contents of the build to the "dapp" directory under the build/ci directory of the system during make to prepare for testing
For example, paracross copy to the chain33/build/ci/paracross directory

Run the test by make docker-compose DAPP="dapp" (here dapp is the name of each dapp such as paracross)
The system will generate a temporary dapp-ci directory based on the build/ci/dapp directory, copy the dapp and system bin and configuration information together to run the test
After the test is executed, the corresponding dockers will not be deleted automatically. You can also view the information through docker exec build_chain33_1 /root/chain33-cli [cmd...]
Use make docker-compose-down DAPP="dapp" to delete the test resources of this dapp and release docker

You can also run all dapps through the dapp parameter keyword all, the all mode will automatically delete the resources of the pass dapp

 1. make docker-compose [proj=xx] [dapp=xx]
    1. If proj and dapp are not set such as make docker-compose, only the test case of the system will be run, and no dapp will be run
    1. If proj is not set, the system will use the build keyword as the service project name of docker-compose by default. If set, the setting shall prevail.
       Different proj can realize docker compose parallel
    1. If dapp is not set, no dapp will be run. If set, only the specified dapp will be run. After the run is over, you need to manually make docker-compose-down dapp=xx release
    1. If dapp=all or ALL, run all dapps that provide testcase
 1. make docker-compose down [proj=xx] [dapp=xx]
    Responsible for clean make docker-compose or make fork-test created docker resources, proj and dapp rules are the same as above
 1. make fork-test [proj=xx] [dapp=xx] fork test
    1. The rules are the same as make docker-compose


## File information under dapp/cmd/build
The files under build are all related to CI testing
 1. Dockerfile, if the dapp has no changes to the Dockerfile of the system, you don’t need to provide it, and use the system default. If there are changes, use your own, the name can remain the same, and the system will not overwrite it.
    You can also change the Dockerfile for a docker, use your own named Dockerfile-xxx, you need to set it in your own docker-compose yml file
    It should be noted that Dockerfile cannot be inherited, only replaced
 1. The docker-compose yml file is organized to organize the docker service. Chain33 needs to start at least 2 dockers to mine. If you modify the Dockerfile naming, you can
    It is specified here that the docker-compose file can be inherited from the system file, that is, only the modified part can be written in the dapp. If there is an overlap with the system, the dapp shall prevail.
    The docker-compose yml file can make various customizations to the docker service
    If there is no modification to docker-compose, it is not necessary to provide
    If used in combination with the system docker-compose.yml file, the compose file of the dapp must conform to the naming of docker-compose-$dapp.yml
 1. testcase.sh The test case is written in this file, and the file name must be testcase.sh, otherwise the system cannot find it.
    Three steps are currently provided for testcase:
    1. init is the modification required to the configuration file before docker starts
    1. config: After docker starts, do the required configuration for dapp, such as transfer, set up wallet, etc.
    1. test: ci test case part
    Testcase must provide a function named after its own dapp as the entry function, three steps are not necessary
    ```
     function paracross() {
         if ["${2}" == "init" ]; then
             para_init
         elif ["${2}" == "config" ]; then
             para_transfer
             para_set_wallet
         elif ["${2}" == "test" ]; then
             para_test
         fi
     
     }
     ```
 1. Fork-test.sh fork test case, it must be this file name, otherwise the system can't find it
    The entrance of fork-test dapp function can be shared with testcase, import testcase.sh through source
    Fork test also provides 5 steps for test configuration
    1. forkInit: Modification of file parameters before docker starts
    1. forkConfig: After the docker service is started, do system configuration, such as transfer
    1. forkAGroupRun: The fork test is divided into two test groups A and B. Alternate mining to simulate the actual fork scenario. Here is the configuration of its own use case when Agroup runs first
    1. ForkBGroupRun: Dapp settings when the B group is running
    1. forkChekRst: Check the results of the system rollback after the fork test is over
    ```
     function privacy() {
         if ["${2}" == "forkInit" ]; then
             privacy_init
         elif ["${2}" == "forkConfig" ]; then
             initPriAccount
         elif ["${2}" == "forkAGroupRun" ]; then
             genFirstChainPritx
             genFirstChainPritxType4
         elif ["${2}" == "forkBGroupRun" ]; then
             genSecondChainPritx
             genSecondChainPritxType4
         elif ["${2}" == "forkCheckRst" ]; then
             checkPriResult
         fi
    
     }
     ```