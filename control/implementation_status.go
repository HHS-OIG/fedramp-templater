package control

import (
	"fmt"
	"github.com/jbowtie/gokogiri/xml"
	"github.com/opencontrol/fedramp-templater/common/implementation"
	"github.com/opencontrol/fedramp-templater/docx"
	"github.com/opencontrol/fedramp-templater/docx/helper"
	"gopkg.in/fatih/set.v0"
)

type implementationStatus struct {
	cell    xml.Node
	statuses map[implementation.Key]*docx.CheckBox
}

func (i *implementationStatus) getCheckedImplementationStatuses() *set.Set {
	// find the control origins currently checked in the section
	checkedImplementations := set.New()
	for status, checkbox := range i.statuses {
		if checkbox.IsChecked() {
			checkedImplementations.Add(status)
		}
	}
	return checkedImplementations
}

func detectImplementationStatusKeyFromDoc(textNodes []xml.Node) implementation.Key {
	textField := helper.ConcatTextNodes(textNodes)
	ImplementationStatusMappings := implementation.GetSourceMappings()
	for implementationStatus, implementationStatusMapping := range ImplementationStatusMappings {
		if implementationStatusMapping.IsDocMappingASubstrOf(textField) {
			return implementationStatus
		}
	}
	return implementation.NoStatus
}

func newImplementationStatus(tbl *table) (*implementationStatus, error) {
	// Find the control origination row.
	rows, err := tbl.Root.Search(".//w:tc[starts-with(normalize-space(.), 'Implementation Status')]")
	if err != nil {
		return nil, err
	}
	// Check that we only found the one cell.
	if len(rows) != 1 {
		return nil, fmt.Errorf("Unable to find Implementation Status cell")
	}
	// Each checkbox is contained in a paragraph.
	statuses := make(map[implementation.Key]*docx.CheckBox)
	paragraphs, err := rows[0].Search(".//w:p")
	if err != nil {
		return nil, err
	}
	for _, paragraph := range paragraphs {
		// 1. Find the box of the checkbox.
		checkBox, err := docx.FindCheckBoxTag(paragraph)
		if err != nil {
			continue
		}

		// 2. Find the text next to the checkbox.
		textNodes, err := paragraph.Search(".//w:t")
		if len(textNodes) < 1 || err != nil {
			continue
		}

		// 3. Detect the key for the map.
		implementationStatusKey := detectImplementationStatusKeyFromDoc(textNodes)
		// if couldn't detect an implementation, skip.
		if implementationStatusKey == implementation.NoStatus {
			continue
		}
		// if the implementation is already in the map, skip.
		_, exists := statuses[implementationStatusKey]
		if exists {
			continue
		}

		// Only construct the checkbox struct if the box and text are found.
		origins[implementationStatusKey] = docx.NewCheckBox(checkBox, &textNodes)
	}
	return &implementationStatus{cell: rows[0], statuses: statuses}, nil
}
