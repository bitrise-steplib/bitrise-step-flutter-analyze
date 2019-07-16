package main

import (
	"testing"
)

func TestHasAnalyzeError(t *testing.T) {
	cases := []struct {
		title     string
		cmdOutput string
		want      bool
	}{
		{
			"no error level violation",
			`info • This class inherits from a class marked as @immutable, and therefore should be immutable (all instance fields must be final) • lib/pages/home/home_page.dart:17:7 • must_be_immutable
			info • This class inherits from a class marked as @immutable, and therefore should be immutable (all instance fields must be final) • lib/pages/login/login_page.dart:11:7 • must_be_immutable
			info • Unused import: 'package:silkthread/service/ads_manager.dart' • lib/pages/post/image_edit.dart:8:8 • unused_import
			info • Unused import: 'dart:ui' • lib/pages/post/image_full_view.dart:1:8 • unused_import
			info • Unused import: 'package:cloud_firestore/cloud_firestore.dart' • lib/pages/post/image_full_view.dart:4:8 • unused_import
			info • Unused import: 'package:font_awesome_flutter/font_awesome_flutter.dart' • lib/pages/post/image_full_view.dart:6:8 • unused_import
			`,
			false,
		},
		{
			"contains error level violation",
			`info • Unused import: 'dart:math' • lib/package.dart:3:8 • unused_import
			error • Expected to find ';' • lib/package.dart:3:8 • expected_token
			`,
			true,
		},
	}

	for _, tt := range cases {
		if got := hasAnalyzeError(tt.cmdOutput); got != tt.want {
			t.Errorf("%s: got %t want %t", tt.title, got, tt.want)
		}
	}
}
