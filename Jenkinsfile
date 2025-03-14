// 브랜치별 배포 위치
def getDeployTargets(envName) {
  targets = [:]

  // dev 브랜치
  targets['master'] = [[
    COMPOSE_ENV: 'master',
    SSH_MODE: 'KEYONLY',
    SSH_KEY_ID: 'toon-server-ssh-key',
    COPY_DIR: '/mnt/blockstorage/containers/blog_back'
  ]]

  return targets[envName]
}

// 브랜치별 환경 정보
def getBuildBranch(branchName) {
  branches = [
    'origin/master': 'master'
  ]

  return branches[branchName]
}

pipeline {
  agent any

  //환경변수
  environment {
    APP_ENV_ID = 'blog-back-env'
    DB_ENV_ID = 'blog-db-env'
    REDIS_ENV_ID = 'blog-redis-env'

    //브랜치 선택
    BUILD_BRANCH = getBuildBranch(env.GIT_BRANCH)

    // 도커 설정
    DOCKER_IMAGE = ''
    RASP_DOCKER_IMAGE = ""

    // DOCKER_IMAGE_NAME = 'stats_service/go_back'
    DOCKER_IMAGE_NAME = 'sjc.vultrcr.com/dongregistry/blog_back'

    SERVER_TARGET = 'toon-server-ip'

    // Git, Docker 레지스트리(https://registry.zetra.kr) 로그인 정보 설정
    GIT_KEY_ID = '2'
    REGISTRY_LOGIN_INFO_ID = 'vultr-registry-user'
  }

  stages {
    stage('체크아웃') {
      steps {
        echo "작업 브랜치: ${env.GIT_BRANCH}"

        git branch: BUILD_BRANCH, credentialsId: GIT_KEY_ID, url: 'git@github.com:Dong-Ian/blog_back_go.git'
      }
    }

  stage('도커 이미지 빌드') {
      steps {
        script {
          DOCKER_IMAGE = docker.build(DOCKER_IMAGE_NAME)
        }

        echo "Built: ${DOCKER_IMAGE_NAME}"
      }
    }

    stage('Raspberry PI 용 도커 이미지 빌드') {
      steps {
        script {
          // 브랜치에 따라 이미지 이름 변경
          RASP_DOCKER_IMAGE = docker.build("${DOCKER_IMAGE_NAME}-raspberry", "--platform=linux/arm64/v8 -f Dockerfile-raspberry .")
        }

        echo "Built: ${RASP_DOCKER_IMAGE}"
      }
    }

    stage('도커 이미지 Push') {
      steps {
        script {
          // 개발서버 내부 Docker 레지스트리(https://registry.zetra.kr)에 업로드
          docker.withRegistry('https://sjc.vultrcr.com/dongregistry', REGISTRY_LOGIN_INFO_ID) {
          // docker.withRegistry('https://registry.donghyuns.com', REGISTRY_LOGIN_INFO_ID) {

            DOCKER_IMAGE.push(env.BUILD_NUMBER)
            DOCKER_IMAGE.push('latest')
            
            RASP_DOCKER_IMAGE.push(env.BUILD_NUMBER)
            RASP_DOCKER_IMAGE.push('latest')
          }
        }

        echo "Pushed: ${DOCKER_IMAGE_NAME}:${env.BUILD_NUMBER}"
      }
    }

  
    stage('도커 컨테이너 배포') {
      steps {
        script {
          def deployTargets = getDeployTargets(BUILD_BRANCH)
          def deployments = [:]

          // 배포 타깃별로 병렬 배포
          for (item in deployTargets) {
            def target = item

            deployments["${BUILD_BRANCH}"] = {
              def remote = [:]

              remote.allowAnyHosts = true

              if (target.SSH_MODE == 'KEYONLY') {
                echo "APP ENV ID: ${APP_ENV_ID}"
                echo "SSH KEY ID: ${target.SSH_KEY_ID}"

                withCredentials([
                  // DOTENV 파일과 SSH KEY를 가져옴
                  file(credentialsId: APP_ENV_ID, variable: 'DOTENV'),
                  file(credentialsId: DB_ENV_ID, variable: 'DB_ENV'),
                  file(credentialsId: REDIS_ENV_ID, variable: 'REDIS_ENV'),

                  string(credentialsId: SERVER_TARGET, variable: 'SSH_IP'),

                  sshUserPrivateKey(credentialsId: target.SSH_KEY_ID, keyFileVariable: 'SSH_PRIVATE_KEY', usernameVariable: 'USERNAME')
                  ]) {
                  remote.name = SSH_IP
                  remote.host = SSH_IP
                  echo "DOTENV file: ${DOTENV}"
                  echo "UserName: ${USERNAME}"
                  echo "COPY DIR: ${SSH_IP}"

                  // 가져온 키로 ssh 정보 설정
                  remote.user = USERNAME
                  remote.identityFile = SSH_PRIVATE_KEY

                  sshCommand remote: remote, command: """
                    mkdir -p ${target.COPY_DIR}/
                  """

                  // docker compose 파일 전송
                  sshPut remote: remote, from: "docker-compose.yml", into:  "${target.COPY_DIR}/docker-compose.yml", failOnError: 'true'

                  // 각 상황에 맞는 .env.* 파일 전송
                  sshPut remote: remote, from: DOTENV, into: "${target.COPY_DIR}/.env", failOnError: 'true'
                  sshPut remote: remote, from: DB_ENV, into: "${target.COPY_DIR}/.db.env", failOnError: 'true'
                  sshPut remote: remote, from: REDIS_ENV, into: "${target.COPY_DIR}/.redis.env", failOnError: 'true'

                  // 도커 이미지 Pull 및 재시작
                  sshCommand remote: remote, command: """
                    cd ${target.COPY_DIR}/
                    sudo BUILD_BRANCH=${BUILD_BRANCH} docker compose pull
                    sudo BUILD_BRANCH=${BUILD_BRANCH} docker compose up -d
                  """
                }
              }
            }

            parallel deployments
          }
        }
      }
    }
  }
}
