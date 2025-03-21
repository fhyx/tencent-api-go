package models

import (
	"testing"
)

func TestTextByName(t *testing.T) {
	tests := []struct {
		name           string
		attributes     Attributes
		searchName     string
		expectedResult string
		expectedFound  bool
	}{
		{
			name: "Found Text Value",
			attributes: Attributes{
				Attrs: []Attribute{
					{Name: "test1", Text: &attrText{Value: "Hello, World!"}},
				},
			},
			searchName:     "test1",
			expectedResult: "Hello, World!",
			expectedFound:  true,
		},
		{
			name: "Found Web URL",
			attributes: Attributes{
				Attrs: []Attribute{
					{Name: "test2", Web: &attrWeb{URL: "https://example.com"}},
				},
			},
			searchName:     "test2",
			expectedResult: "https://example.com",
			expectedFound:  true,
		},
		{
			name: "No Matching Name",
			attributes: Attributes{
				Attrs: []Attribute{
					{Name: "test1", Text: &attrText{Value: "Hello, World!"}},
				},
			},
			searchName:     "nonexistent",
			expectedResult: "",
			expectedFound:  false,
		},
		{
			name: "Attribute Without Text or Web",
			attributes: Attributes{
				Attrs: []Attribute{
					{Name: "test3", MiniApp: &attrMiniApp{AppID: "12345"}},
				},
			},
			searchName:     "test3",
			expectedResult: "",
			expectedFound:  false,
		},
		{
			name: "Multiple Attributes, Find Correct One",
			attributes: Attributes{
				Attrs: []Attribute{
					{Name: "test1", Text: &attrText{Value: "Hello, World!"}},
					{Name: "test2", Web: &attrWeb{URL: "https://example.com"}},
				},
			},
			searchName:     "test2",
			expectedResult: "https://example.com",
			expectedFound:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, found := tt.attributes.TextByName(tt.searchName)
			if result != tt.expectedResult || found != tt.expectedFound {
				t.Errorf("TextWithName(%q) = (%q, %v), expected (%q, %v)", tt.searchName, result, found, tt.expectedResult, tt.expectedFound)
			}
		})
	}
}
