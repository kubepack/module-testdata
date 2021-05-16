package action

import (
	"fmt"
	"log"

	"helm.sh/helm/v3/pkg/chart"
)

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	_ = log.Output(2, fmt.Sprintf(format, v...))
}

func setAnnotations(chrt *chart.Chart, k, v string) {
	if chrt.Metadata.Annotations == nil {
		chrt.Metadata.Annotations = map[string]string{}
	}
	if k != "" {
		chrt.Metadata.Annotations[k] = v
	} else {
		delete(chrt.Metadata.Annotations, k)
	}
}
