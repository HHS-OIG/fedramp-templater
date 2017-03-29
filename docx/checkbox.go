package docx

import (
	"fmt"
	"strconv"
	"github.com/jbowtie/gokogiri/xml"
	"github.com/opencontrol/fedramp-templater/docx/helper"
)

const (
	checkBoxAttributeKey    = "val"
	checkBoxCheckedValue    = "1"
	checkBoxNotCheckedValue = "0"
	checkBoxCheckedChild	= "checked"
	checkBoxCheckedStateChild	= "checkedState"
	checkBoxUncheckedStateChild	= "uncheckedState"
	checkBoxUnicodeKey		= "w:sdtContent"
)

// NewCheckBox constructs a new checkbox. Checks if the checkmark value can actually be found.
// If it cannot be found, will return nil.
func NewCheckBox(checkMark xml.Node, textNodes *[]xml.Node) *CheckBox {
	// Have to use Attr.
	// Using Attribute does not work when checking the value key.
	// Make sure the length is non zero.
	checkedNode, err := getChild(checkMark, checkBoxCheckedChild)
	if err != nil {
		return nil
	}
	if len(checkedNode.Attr(checkBoxAttributeKey)) == 0 {
		return nil
	}
	return &CheckBox{checkMark: checkMark, textNodes: textNodes}
}

// CheckBox represents a checkbox in a word document with any corresponding text.
type CheckBox struct {
	checkMark xml.Node
	textNodes *[]xml.Node
}

func getChild(c xml.Node, childName string) (xml.Node, error) {
    if c == nil {
		return nil, fmt.Errorf("Node provided has value 'nil'")
	}

    searchQuery := fmt.Sprintf(".//*[local-name()='%s']", childName)
	children, err := c.Search(searchQuery)
	if err != nil {
		return nil, err
	} else if len(children) != 1 {
		return nil, fmt.Errorf("Unable to find the check box child.")
	}

	return children[0], nil
}

// IsChecked will return true if the box is checked, false otherwise.
func (c *CheckBox) IsChecked() bool {
	child, err := getChild(c.checkMark, checkBoxCheckedChild)
    if err != nil {
		return false
	}
	return child.Attr(checkBoxAttributeKey) == checkBoxCheckedValue
}

// SetCheckMarkTo will set the checkbox state according to the input value.
func (c *CheckBox) SetCheckMarkTo(value bool) {
    // First, set checkBoxCheckedChild
	checkBoxValue := checkBoxNotCheckedValue
	if value == true {
		checkBoxValue = checkBoxCheckedValue
	}
	checkedNode, err := getChild(c.checkMark, checkBoxCheckedChild)
	if err != nil {
		return
	}
	checkedNode.AttributeList()[0].SetContent(checkBoxValue)

	// Next, grab the appropriate unicode value to insert
	var unicodeKey string
	if value == true {
		unicodeKey = checkBoxCheckedStateChild
	} else {
		unicodeKey = checkBoxUncheckedStateChild
	}
	unicodeKeyNode, err := getChild(c.checkMark, unicodeKey)
	if err != nil {
		return
	}
    searchQuery := fmt.Sprintf(".//%s//w:t", checkBoxUnicodeKey)
	unicodeNodes, err := c.checkMark.Parent().Parent().Search(searchQuery)
	if err != nil {
		return
	} else if len(unicodeNodes) != 1 {
		return
	}
	// Convert hex value (in string format) to proper integer value
	unicodeKeyIntValue, err := strconv.ParseInt(unicodeKeyNode.Attr(checkBoxAttributeKey), 16, 32)
	if err != nil {
		return
	}
	unicodeNodes[0].SetContent(string(unicodeKeyIntValue))
}

// GetTextValue will return the corresponding text for the checkbox.
func (c *CheckBox) GetTextValue() string {
	return helper.ConcatTextNodes(*(c.textNodes))
}

// FindCheckBoxTag will look for a checkbox and the value tag.
// We look for the w:default tag embedded in the w:checkBox tag because that is what we need to modify the checkbox.
func FindCheckBoxTag(paragraph xml.Node) (xml.Node, error) {
	checkBoxes, err := paragraph.Search(".//*[local-name()='checkbox']")
	if err != nil {
		return nil, err
	} else if len(checkBoxes) != 1 {
		return nil, fmt.Errorf("Unable to find the check box value.")
	}
	return checkBoxes[0], nil
}
