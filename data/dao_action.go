package data

import model "github.com/awlsring/action-runner-model"

type Action interface {
	Run(e model.RunActionResponseContent) error
}