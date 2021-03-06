pipeline {
    agent any
    environment {
        PATH = "/usr/local/go/bin:$PATH"
        DOCKER_IMAGE_NAME = "harbor.vmware.local/library/go-cicd-kubernetes"
    }
    stages {
        stage('Build') {
            steps {
                echo 'Compiling Program'
                sh 'go test -v'
            }
        }
        stage('Build Docker Image') {
            when { 
                branch 'master'
            }
            steps {
                script {
                    app = docker.build(DOCKER_IMAGE_NAME)
                    app.withRun("-d -p 8181:8181") { c ->
                        sh 'curl localhost:8181'
                    }    
                }
            }
        }
        stage('Push Docker Image') {
            when {
                branch 'master'
            }
            steps {
                script {
                    docker.withRegistry('https://harbor.vmware.local', 'harbor') {
                        app.push("${env.BUILD_NUMBER}")
                        app.push("latest")
                    }
                }
            }
        }
        stage('DeployToPKS') {
            when {
                branch 'master'
            }
            steps {
                milestone(1)
                withCredentials([usernamePassword(credentialsId: 'pks_client', usernameVariable: 'USERNAME', passwordVariable: 'USERPASS')]) {
                    script {
                        //sh "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$pks_client \"kubectl delete all -l app=gocicd\""
                        if (sh(script: "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$pks_client \"kubectl get deployment gocicd\"" , returnStatus: true) == 0) {
                            echo "Deployment exists"
                            sh "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$pks_client \"kubectl set image deployment gocicd gocicd=${DOCKER_IMAGE_NAME}:${BUILD_NUMBER}\""
                        } else {
                            echo "Deployment doesn't exist"
                            sh "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$pks_client \"kubectl run gocicd --image ${DOCKER_IMAGE_NAME}:${BUILD_NUMBER}\""
                        }
			retry (10) {
			    if (sh(script: "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$pks_client \"kubectl get svc gocicd\"" , returnStatus: true) == 0) {
                                echo "Service exists"
                            } else {
                                echo "Service doesn't exist"
                                sh "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$pks_client \"kubectl expose deployment gocicd --type LoadBalancer --port 80 --target-port 8181\""
                            }
			}
                    }
                }
                /*
                kubernetesDeploy(
                  kubeconfigId: 'kubeconfig',
                  configs: 'kubernetes.yaml',
                  enableConfigSubstitution: true
                ) 
                */
            }
        }
        stage('Get Service IP') {
            when {
                branch 'master'
            }
            steps {
	        /*
                script {
                    def result = sh(script: "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$pks_client \"ls\"" , returnStdout: true)     
                }
		*/
                //milestone(1)
                
		retry(10) {
                    withCredentials([usernamePassword(credentialsId: 'pks_client', usernameVariable: 'USERNAME', passwordVariable: 'USERPASS')]) {
                        script {
                            //sh "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$pks_client \"kubectl get all -l app=gocicd -o wide\""
                            def ip = sh (script: "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$pks_client \"kubectl get svc gocicd --output=jsonpath={'.status.loadBalancer.ingress[].ip'}\"", returnStdout: true)
                            sh 'sleep 5'
                            echo "IP is ${ip}"
                            echo "URL is http://${ip}:8181"
                            try {
                             //echo sh(script: "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$kube-master-ip ls -la", returnStdout: true)
                             //sh "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$kube-master-ip \"/usr/bin/touch /tmp/jenkins\""
                            } catch (err) {
                             echo: 'caught error: $err'
                            }
                        }
                    }
                }
            }
        }
    }
}
