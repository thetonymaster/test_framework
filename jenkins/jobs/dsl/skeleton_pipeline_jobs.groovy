
import pluggable.scm.*

@Grab(group='org.yaml', module='snakeyaml', version='1.18') 
import org.yaml.snakeyaml.Yaml

// Folders
def workspaceFolderName = "${WORKSPACE_NAME}"
def projectFolderName = "${PROJECT_NAME}"

Yaml yaml = new Yaml()

def variables = yaml.load(("/var/jenkins_home/jobs/test/jobs/test_project/jobs/LoadDevCartridge/workspace/cartridge/jenkins/jobs/dsl" + "/example.yaml" as File).text)

// Variables
// **The git repo variables will be changed to the users' git repositories manually in the Jenkins jobs**
def referenceAppgitRepo = "spring-petclinic"
def referenceAppGitUrl = "ssh://jenkins@gerrit:29418/${PROJECT_NAME}/" + referenceAppgitRepo

// def regressionTestGitRepo = "YOUR_REGRESSION_TEST_REPO"


// ** The logrotator variables should be changed to meet your build archive requirements
def logRotatorDaysToKeep = 7
def logRotatorBuildNumToKeep = 7
def logRotatorArtifactsNumDaysToKeep = 7
def logRotatorArtifactsNumToKeep = 7

// Jobs
def buildAppJob = freeStyleJob(projectFolderName + "/Test_Framework_Build")
def analyzeScriptsJob = freeStyleJob(projectFolderName + "/Test_Framework_Analyze_Scripts")
def analyzeTestsJob = freeStyleJob(projectFolderName + "/Test_Framework_Analyze_Tests")

// def unitTestJob = freeStyleJob(projectFolderName + "/Skeleton_Application_Unit_Tests")
// def codeAnalysisJob = freeStyleJob(projectFolderName + "/Skeleton_Application_Code_Analysis")
// def deployJob = freeStyleJob(projectFolderName + "/Skeleton_Application_Deploy")
// def regressionTestJob = freeStyleJob(projectFolderName + "/Skeleton_Application_Regression_Tests")

// Views
def pipelineView = buildPipelineView(projectFolderName + "/Test_Framework_Application")

pipelineView.with{
    title('Skeleton Application Pipeline')
    displayedBuilds(5)
    selectedJob(projectFolderName + "/Test_Framework_Build")
    showPipelineParameters()
    showPipelineDefinitionHeader()
    refreshFrequency(5)
}




// All jobs are tied to build on the Jenkins slave
// The functional build steps for each job have been left empty
// A default set of wrappers have been used for each job
// New jobs can be introduced into the pipeline as required

buildAppJob.with{
  description("Build job ")
  logRotator {
    daysToKeep(logRotatorDaysToKeep)
    numToKeep(logRotatorBuildNumToKeep)
    artifactDaysToKeep(logRotatorArtifactsNumDaysToKeep)
    artifactNumToKeep(logRotatorArtifactsNumToKeep)
  }
  
  scm {
        git {
            remote {
                url(referenceAppGitUrl)
                credentials("adop-jenkins-master")
            }
            branch("*/master")
        }
    }

  environmentVariables {
      env('WORKSPACE_NAME',workspaceFolderName)
      env('PROJECT_NAME',projectFolderName)
  }
  label("java8")
  triggers {
      gerrit {
          events {
              refUpdated()
          }
          project(projectFolderName + '/' + referenceAppgitRepo, 'plain:master')
          configure { node ->
              node / serverName("ADOP Gerrit")
          }
      }
  }
  steps {
      maven {
          goals('clean install -DskipTests')
          mavenInstallation("ADOP Maven")
      }
  }
  publishers{
    downstreamParameterized{
      trigger(projectFolderName + "/Test_Framework_Analyze_Scripts"){
        condition("UNSTABLE_OR_BETTER")
        parameters{
          predefinedProp("B",'${BUILD_NUMBER}')
          predefinedProp("PARENT_BUILD", '${JOB_NAME}')
        }
      }
    }
  }
}

analyzeScriptsJob.with{
  description("This job analyzes the test scripts.")
  parameters{
    stringParam("B",'',"Parent build number")
    stringParam("PARENT_BUILD","Test_Framework_Build","Parent build name")
  }
  logRotator {
    daysToKeep(logRotatorDaysToKeep)
    numToKeep(logRotatorBuildNumToKeep)
    artifactDaysToKeep(logRotatorArtifactsNumDaysToKeep)
    artifactNumToKeep(logRotatorArtifactsNumToKeep)
  }
  environmentVariables {
      env('WORKSPACE_NAME',workspaceFolderName)
      env('PROJECT_NAME',projectFolderName)
  }
  wrappers {
    preBuildCleanup()
    injectPasswords()
    maskPasswords()
    sshAgent("adop-jenkins-master")
  }
  label("docker")
  steps {
    shell('''ls -lah'''.stripMargin())
  }
  publishers{
    downstreamParameterized{
      trigger(projectFolderName + "/Test_Framework_Analyze_Tests"){
        condition("UNSTABLE_OR_BETTER")
        parameters{
          predefinedProp("B",'${B}')
          predefinedProp("PARENT_BUILD", '${PARENT_BUILD}')
        }
      }
    }
  }
}

analyzeTestsJob.with{
  description("This job analyzes the tests.")
  parameters{
    stringParam("B",'',"Parent build number")
    stringParam("PARENT_BUILD","Test_Framework_Build","Parent build name")
  }
  logRotator {
    daysToKeep(logRotatorDaysToKeep)
    numToKeep(logRotatorBuildNumToKeep)
    artifactDaysToKeep(logRotatorArtifactsNumDaysToKeep)
    artifactNumToKeep(logRotatorArtifactsNumToKeep)
  }
  environmentVariables {
      env('WORKSPACE_NAME',workspaceFolderName)
      env('PROJECT_NAME',projectFolderName)
  }
  wrappers {
    preBuildCleanup()
    injectPasswords()
    maskPasswords()
    sshAgent("adop-jenkins-master")
  }
  label("docker")
  steps {
    shell('''## YOUR CODE ANALYSIS STEPS GO HERE'''.stripMargin())
  }
  publishers{
    downstreamParameterized{
      trigger(projectFolderName + "/Test_Framework_Analyze_Scripts"){
        condition("UNSTABLE_OR_BETTER")
        parameters{
          predefinedProp("B",'${B}')
          predefinedProp("PARENT_BUILD", '${PARENT_BUILD}')
        }
      }
    }
  }
}

// deployJob.with{
//   description("This job deploys the skeleton application to the CI environment")
//   parameters{
//     stringParam("B",'',"Parent build number")
//     stringParam("PARENT_BUILD","Skeleton_Application_Build","Parent build name")
//     stringParam("ENVIRONMENT_NAME","CI","Name of the environment.")
//   }
//   logRotator {
//     daysToKeep(logRotatorDaysToKeep)
//     numToKeep(logRotatorBuildNumToKeep)
//     artifactDaysToKeep(logRotatorArtifactsNumDaysToKeep)
//     artifactNumToKeep(logRotatorArtifactsNumToKeep)
//   }
//   wrappers {
//     preBuildCleanup()
//     injectPasswords()
//     maskPasswords()
//     sshAgent("adop-jenkins-master")
//   }
//   environmentVariables {
//       env('WORKSPACE_NAME',workspaceFolderName)
//       env('PROJECT_NAME',projectFolderName)
//   }
//   label("docker")
//   steps {
//     shell('''## YOUR DEPLOY STEPS GO HERE'''.stripMargin())
//   }
//   publishers{
//     downstreamParameterized{
//       trigger(projectFolderName + "/Skeleton_Application_Regression_Tests"){
//         condition("UNSTABLE_OR_BETTER")
//         parameters{
//           predefinedProp("B",'${B}')
//           predefinedProp("PARENT_BUILD", '${PARENT_BUILD}')
//           predefinedProp("ENVIRONMENT_NAME", '${ENVIRONMENT_NAME}')
//         }
//       }
//     }
//   }
// }
// 
// regressionTestJob.with{
//   description("This job runs regression tests on the deployed skeleton application")
//   parameters{
//     stringParam("B",'',"Parent build number")
//     stringParam("PARENT_BUILD","Skeleton_Application_Build","Parent build name")
//     stringParam("ENVIRONMENT_NAME","CI","Name of the environment.")
//   }
//   logRotator {
//     daysToKeep(logRotatorDaysToKeep)
//     numToKeep(logRotatorBuildNumToKeep)
//     artifactDaysToKeep(logRotatorArtifactsNumDaysToKeep)
//     artifactNumToKeep(logRotatorArtifactsNumToKeep)
//   }
//   scm scmProvider.get(projectFolderName, regressionTestGitRepo, "*/master", "adop-jenkins-master", null)
//   wrappers {
//     preBuildCleanup()
//     injectPasswords()
//     maskPasswords()
//     sshAgent("adop-jenkins-master")
//   }
//   environmentVariables {
//       env('WORKSPACE_NAME',workspaceFolderName)
//       env('PROJECT_NAME',projectFolderName)
//   }
//   label("docker")
//   steps {
//     shell('''## YOUR REGRESSION TESTING STEPS GO HERE'''.stripMargin())
//   }
// }
