pipeline {
    agent any

    environment {
        IMAGE_NAME = "tilalx/meme-api"
        DOCKER_BUILDKIT = 1
        PIPELINE_NAME = "${JOB_NAME.replaceAll('/', '_')}-${BUILD_NUMBER}"
        DOCKER_CLI_EXPERIMENTAL = 'enabled'
    }

    options {
        timestamps()
        disableConcurrentBuilds()
    }

    stages {
        stage('Setup Buildx') {
            steps {
                script {
                    sh 'docker run --rm --privileged tonistiigi/binfmt --install all'
                    def builderName = "builder-${env.BUILD_ID}-${env.BRANCH_NAME}"
                    sh "docker buildx create --name ${builderName} --use"
                    sh 'docker buildx inspect --bootstrap'
                    withCredentials([usernamePassword(credentialsId: 'dockerhub', usernameVariable: 'DOCKERHUB_USER', passwordVariable: 'DOCKERHUB_PASS')]) {
                        sh 'echo $DOCKERHUB_PASS | docker login -u $DOCKERHUB_USER --password-stdin'
                    }
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    def branchName = env.BRANCH_NAME
                    def tagName = ""
                    def isReleaseCommit = false
                    def releaseVersion = ""

                    // Retrieve and convert the latest commit message to lowercase
                    def commitMessage = sh(script: 'git log -1 --pretty=%B', returnStdout: true).trim().toLowerCase()

                    // Check if the commit message contains "release x.x.x"
                    if (commitMessage ==~ /.*release\s+(\d+\.\d+\.\d+).*/) {
                        def matcher = (commitMessage =~ /release\s+(\d+\.\d+\.\d+)/)
                        if (matcher) {
                            releaseVersion = matcher[0][1]
                            isReleaseCommit = true
                        }
                    }

                    if (branchName == "main") {
                        tagName = "rolling"
                    } else if (branchName.startsWith("PR-")) {
                        tagName = "pr-${branchName.split('-')[1]}"
                    } else {
                        tagName = branchName
                    }

                    def tags = "-t ${env.IMAGE_NAME}:${tagName}"

                    if (isReleaseCommit) {
                        tags += " -t ${env.IMAGE_NAME}:latest -t ${env.IMAGE_NAME}:${releaseVersion}"
                    }

                    sh """
                        docker buildx build --platform linux/amd64 --provenance=true --sbom=true --build-arg DOCKER_BUILDKIT=${DOCKER_BUILDKIT} --memory 32g --memory-swap 16g ${tags} --push .
                    """
                }
            }
        }
    }

    post {
        always {
            script {
                def builderName = "builder-${env.BUILD_ID}-${env.BRANCH_NAME}"
                sh "docker buildx rm ${builderName}"
            }
            cleanWs()
        }
    }
}
