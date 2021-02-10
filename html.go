/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	// standard packages
	"bytes"
	"regexp"

	// external packages
	"golang.org/x/net/html"
)

func getDocumentTree(source []byte) (*html.Node, error) {
	return html.Parse(bytes.NewReader(source))
}

func getElementByID(node *html.Node, idPat string) *html.Node {
	return getElementByIDAndTag(node, idPat, "")
}

func getElementByTag(node *html.Node, tag string) *html.Node {
	return getElementByIDAndTag(node, "", tag)
}

func getElementByIDAndTag(node *html.Node, idPat, tag string) *html.Node {
	if node == nil {
		return nil
	}

	var (
		tagOK = false
		idOK  = false
	)
	if node.Type == html.ElementNode {
		if tag == "" || tag == node.Data {
			tagOK = true
		}
		if idPat == "" {
			idOK = true
		} else if len(node.Attr) > 0 {
			re := regexp.MustCompile(idPat)
			for _, attr := range node.Attr {
				if attr.Key == "id" && re.MatchString(attr.Val) {
					idOK = true
				}
			}
		}
		if tagOK && idOK {
			return node
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if result := getElementByIDAndTag(child, idPat, tag); result != nil {
			return result
		}
	}
	return nil
}

func hasAttr(node *html.Node, attrName string) bool {
	for _, attr := range node.Attr {
		if attrName == attr.Key {
			return true
		}
	}
	return false
}

func hasAttrRe(node *html.Node, attrRe *regexp.Regexp) bool {
	for _, attr := range node.Attr {
		if attrRe.MatchString(attr.Key) {
			return true
		}
	}
	return false
}

func getAttr(node *html.Node, attrName string) *html.Attribute {
	for _, attr := range node.Attr {
		if attrName == attr.Key {
			return &attr
		}
	}
	return nil
}

func getAttrRe(node *html.Node, attrRe *regexp.Regexp) *html.Attribute {
	for _, attr := range node.Attr {
		if attrRe.MatchString(attr.Key) {
			return &attr
		}
	}
	return nil
}
