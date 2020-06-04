// a monolithic jenkinsfile
// deal with payload sent by gitlab
// resolve gitlab action type among push, merge request, tag
// referenced: https://github.com/jenkinsci/gitlab-plugin
pipeline {
    agent any
    tools {
        go 'Go 1.14'
    }

    environment {
        PACK_FOLDER = "./chaos_tmp_pack"
    }

    stages {
        stage('Checkout') {
            steps {
                script {
                    try {
                        switch(gitlabActionType) {
                        // case "MERGE":
                        //     echo "repo: ${gitlabSourceRepoHomepage} user: ${gitlabUserName} email:${gitlabUserEmail} action: ${gitlabActionType} source ${gitlabSourceBranch} target ${gitlabTargetBranch}"
                        //     // sh "git checkout ${gitlabSourceBranch}"
                        //     sh "git checkout ${gitlabTargetBranch}"
                        //     sh "git merge origin/${gitlabSourceBranch} -m 'merge'"
                        //     break
                        case "PUSH":
                            // echo "repo: ${gitlabSourceRepoHomepage} user: ${gitlabUserName} email:${gitlabUserEmail} action: ${gitlabActionType} before: ${gitlabBefore} after: ${gitlabAfter}"
                            sh "git checkout ${gitlabAfter}"
                            break
                        case "TAG_PUSH":
                            // echo "repo: ${gitlabSourceRepoHomepage} user: ${gitlabUserName} email:${gitlabUserEmail} action: ${gitlabActionType} before: ${gitlabBefore} after: ${gitlabAfter}"
                            sh "git checkout ${gitlabAfter}"
                            break
                        default:
                            echo gitlabActionType
                        }
                    } catch (Exception ex){
                        echo "push by hand"
                    }
                }
            }
        }

        stage('Build') {
            steps {
                script {
                    sh 'go version'
                    sh 'rm -rf ./oxygen.tar.gz'
                    sh 'rm -rf ${PACK_FOLDER}'
                    sh 'mkdir ${PACK_FOLDER}'
                    sh 'cp -rf ./config ${PACK_FOLDER}/'
                    sh 'cp -rf ./script ${PACK_FOLDER}/'
                    sh 'rm -rf ${PACK_FOLDER}/config/db.toml'
                    if (gitlabActionType == 'TAG_PUSH') {
                        sh 'PATH=$PATH:/var/jenkins_home/go/bin GOPROXY=https://goproxy.io; GOSUMDB=off; make prod;'
                    } else {
                        sh 'PATH=$PATH:/var/jenkins_home/go/bin GOPROXY=https://goproxy.io; GOSUMDB=off; make test;'
                    }
                    sh 'mv oxygen ${PACK_FOLDER}/'
                    sh 'tar -czf ./oxygen.tar.gz ${PACK_FOLDER}/*'
                }
            }
        }

        stage('Deploy') {
            steps {
                script {
                    if (gitlabActionType == 'TAG_PUSH') {
                        echo 'TAG_PUSH'
                        NEWEST_TAG = sh(returnStdout: true, script: 'git describe --abbrev=0 --tags').trim()
                        NOWTIME = sh(returnStdout: true, script: "date '+%Y-%m-%d-%H-%M-%S'").trim()
                        FILE_NAME = "chaos_${NEWEST_TAG}_${NOWTIME}.tar.gz"
                        sh "mv oxygen.tar.gz ${FILE_NAME}"
                        sh "scp ${FILE_NAME} admin@1127.0.0.1:/data/oxygen/"
                        sh "mv ${FILE_NAME} ../oxygen/"
                        echo "最新包名:${FILE_NAME}"
                    } else {
                        echo 'MERGE'
                        sh 'ssh root@127.0.0.1 "mkdir -p /data//oxygen"'
                        sh 'scp ./oxygen.tar.gz root@127.0.0.1:/data/tds/outer/'
                        sh 'ssh root@127.0.0.1 "cd /data/tds/outer/; rm -rf ${PACK_FOLDER}; tar -xzf oxygen.tar.gz;"'
                        sh 'ssh root@127.0.0.1 "cd /data/tds/outer/; cp -rf ${PACK_FOLDER}/* oxygen"'
                        sh 'JENKINS_NODE_COOKIE=dontKillMe; ssh root@127.0.0.1 "/data/tds/outer/oxygen/script/restart.sh"'
                    }
                }
            }
        }

        stage('Clean') {
            steps {
                cleanWs()
            }
        }
    }
    post {
        failure {
            dingTalk accessToken: '', jenkinsUrl: 'http://127.0.0.1:10086/',
                    message: "chaos部署失败。", notifyPeople: 'Jenkins'
        }
        // success {
        //     dingTalk accessToken: '', jenkinsUrl: '',
        //             message: "部署成功。", notifyPeople: 'Jenkins'
        // }
    }
}

