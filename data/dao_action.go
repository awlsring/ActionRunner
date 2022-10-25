package data

import model "github.com/awlsring/dws-action-runner"

type Action interface {
	Run(e model.RunActionResponseContent) error
}