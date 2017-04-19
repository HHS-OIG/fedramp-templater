package control

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jbowtie/gokogiri/xml"
	"github.com/opencontrol/fedramp-templater/xml/helper"
	"gopkg.in/fatih/set.v0"
)

type Parameter struct {
	parentNode xml.Node
	textNodes  *[]xml.Node
}

func NewParameter(parentNode xml.Node, textNodes *[]xml.Node) *Parameter {
    return &Parameter{parentNode: parentNode, textNodes: textNodes}
}

// findParameters looks for the Parameter cell(s) in the control table.
func findParameters(ct *SummaryTable) (*set.Set, error) {
	parameterNodeSet := set.New()
	parameterNodes, err := ct.table.searchSubtree(".//w:tc[starts-with(normalize-space(.), 'Parameter')]")
	if (err == nil && len(parameterNodes) >= 1) {
	    for _, node := range parameterNodes {
	        childNodes, childErr := helper.SearchSubtree(node, `.//w:t`)
	        if childErr != nil || len(childNodes) < 1 {
		        return nil, errors.New("Should not happen, cannot find text nodes.")
	        }
			parameterNodeSet.Add(&Parameter{parentNode: node, textNodes: &childNodes})
        }
	}
	return parameterNodeSet, err
}

// parameter is the container for the responsible role cell.
// getContent returns the full string representation of the content of the cell itself.
func (r *Parameter) getContent() string {
	return r.parentNode.Content()
}

// setValue will set the value of the responsible role cell and do any needed formatting.
// In this case, it will just place the text after ":"
// If there are other nodes, we don't care about them, zero the content out.
func (r *Parameter) setValue(value string) {
    id := r.getId()

	for idx, node := range *(r.textNodes) {
		if idx == 0 {
			node.SetContent(fmt.Sprintf("Parameter %s: %s", id, value))
		} else {
			node.SetContent("")
		}
	}
}

// isDefaultValue contains the logic to detect if the input is a default value. This is looking at the extracted
// value of responsible role and not the full string representation.
func (r *Parameter) isDefaultValue(value string) bool {
	return value == ""
}

// getId returns the ID from the full string representation.
// It looks at the text after Parameter and before ":"
func (r *Parameter) getId() string {
	re := regexp.MustCompile("Parameter (.+):(.+)")
	idText := re.FindStringSubmatch(r.parentNode.Content())[1]
	idTextNoSpaces := strings.Replace(idText, " ", "", -1)
	return strings.TrimSpace(idTextNoSpaces)
}

// getValue extracts the unique value from the full string representation.
// It looks at all the text after ":".
func (r *Parameter) getValue() string {
	re := regexp.MustCompile("Parameter (.+):(.+)")
	parameterText := re.FindStringSubmatch(r.parentNode.Content())[2]
	return strings.TrimSpace(parameterText)
}
