package service

import (
	edpv1alpha1 "cd-pipeline-handler-controller/pkg/apis/edp/v1alpha1"
	jenkinsClient "cd-pipeline-handler-controller/pkg/jenkins"
	ClientSet "cd-pipeline-handler-controller/pkg/openshift"
	"context"
	"errors"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type CDPipelineService struct {
	Resource *edpv1alpha1.CDPipeline
	Client   client.Client
}

const (
	StatusInit       = "initialized"
	StatusFailed     = "failed"
	StatusFinished   = "created"
	StatusInProgress = "in progress"
)

func (s CDPipelineService) CreateCDPipeline() error {
	cr := s.Resource

	log.Printf("Start creating CD Pipeline: %v", cr.Spec.Name)
	if cr.Status.Status != StatusInit {
		log.Printf("CD Pipeline %v is not in init status. Skipped", cr.Spec.Name)
		return errors.New(fmt.Sprintf("CD Pipeline %v is not in init status. Skipped", cr.Spec.Name))
	}
	log.Printf("CD Pipeline %v has 'init' status", cr.Spec.Name)

	err := s.updateStatus(edpv1alpha1.CDPipelineStatus{
		Status:          StatusInProgress,
		Available:       false,
		LastTimeUpdated: time.Now(),
		Action:          edpv1alpha1.AcceptCDPipelineRegistration,
		Result:          edpv1alpha1.Success,
		Username:        "system",
		Value:           "inactive",
	})
	if err != nil {
		return fmt.Errorf("error has been occurred in cd_pipeline status update: %v", err)
	}

	jenkinsUrl := fmt.Sprintf("http://jenkins.%s:8080", cr.Namespace)
	log.Printf("Jenkins URL has been generated: %v", jenkinsUrl)

	jenkinsToken, jenkinsUsername, err := getJenkinsCreds(ClientSet.CreateOpenshiftClients(), cr.Namespace)
	if err != nil {
		log.Println("Couldn't fetch Jenkins creds")
		s.setFailedFields(edpv1alpha1.JenkinsConfiguration, err.Error())
		return err
	}

	jenkins, err := jenkinsClient.Init(jenkinsUrl, jenkinsUsername, jenkinsToken)
	if err != nil {
		log.Println("Couldn't initialize Jenkins client")
		s.setFailedFields(edpv1alpha1.JenkinsConfiguration, err.Error())
		return err
	}

	_, err = jenkins.CreateFolder(cr.Name + "-cd-pipeline")
	if err != nil {
		log.Println("Couldn't create folder for Jenkins")
		s.setFailedFields(edpv1alpha1.JenkinsConfiguration, err.Error())
		return err
	}

	err = s.updateStatus(edpv1alpha1.CDPipelineStatus{
		Status:          StatusFinished,
		Available:       true,
		LastTimeUpdated: time.Now(),
		Username:        "system",
		Action:          edpv1alpha1.SetupInitialStructureForCDPipeline,
		Result:          edpv1alpha1.Success,
		Value:           "active",
	})
	if err != nil {
		return fmt.Errorf("error has been occurred in cd_pipeline status update: %v", err)
	}

	log.Printf("CD pipeline has been created. Status: %v", StatusFinished)
	return nil
}

func (s CDPipelineService) updateStatus(status edpv1alpha1.CDPipelineStatus) error {
	s.Resource.Status = status

	err := s.Client.Status().Update(context.TODO(), s.Resource)
	if err != nil {
		err := s.Client.Update(context.TODO(), s.Resource)
		if err != nil {
			return err
		}
	}

	log.Printf("Status for CD Pipeline %v is set up.", s.Resource.Name)

	return nil
}

func (s CDPipelineService) setFailedFields(action edpv1alpha1.ActionType, message string) {
	s.Resource.Status = edpv1alpha1.CDPipelineStatus{
		Status:          StatusFailed,
		Available:       false,
		LastTimeUpdated: time.Now(),
		Username:        "system",
		Action:          action,
		Result:          edpv1alpha1.Error,
		DetailedMessage: message,
		Value:           "failed",
	}

	log.Printf("Status %v for CD Pipeline %v is set up.", edpv1alpha1.Error, s.Resource.Name)
}

func getJenkinsCreds(clientSet *ClientSet.ClientSet, namespace string) (string, string, error) {
	log.Printf("Start recieving credentials for Jenkins in namespace %v", namespace)
	jenkinsTokenSecret, err := clientSet.CoreClient.Secrets(namespace).Get("jenkins-token", metav1.GetOptions{})
	if err != nil {
		errorMsg := fmt.Sprint(err)
		log.Println(errorMsg)
		return "", "", errors.New(errorMsg)
	}

	log.Printf("Credentials for Jenkins in namespace %v has been recieved", namespace)

	return string(jenkinsTokenSecret.Data["token"]), string(jenkinsTokenSecret.Data["username"]), nil
}
