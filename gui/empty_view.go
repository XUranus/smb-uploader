package gui

import "github.com/lxn/walk"

func CheckScrollViewEmpty(scrollView *walk.ScrollView, emptyPanel *walk.Composite)  {
	if scrollView.Children().Len() == 0 {
		emptyPanel.SetVisible(true)
	}
}
