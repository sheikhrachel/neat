// Copyright © 2019 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package neat

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gonvenience/bunt"
	colorful "github.com/lucasb-eyer/go-colorful"
)

// BoxStyle represents a styling option for a content box
type BoxStyle func(*boxOptions)

type boxOptions struct {
	headlineColor  *colorful.Color
	contentColor   *colorful.Color
	headlineStyles []bunt.StyleOption
}

// HeadlineColor sets the color of the headline text
func HeadlineColor(color colorful.Color) BoxStyle {
	return func(options *boxOptions) {
		options.headlineColor = &color
	}
}

// HeadlineStyle sets the style to be used for the headline text
func HeadlineStyle(style bunt.StyleOption) BoxStyle {
	return func(options *boxOptions) {
		options.headlineStyles = append(options.headlineStyles, style)
	}
}

// ContentColor sets the color of the content text
func ContentColor(color colorful.Color) BoxStyle {
	return func(options *boxOptions) {
		options.contentColor = &color
	}
}

// ContentBox creates a string for the terminal where content is printed inside
// a simple box shape.
func ContentBox(headline string, content string, opts ...BoxStyle) string {
	var buf bytes.Buffer
	Box(&buf, headline, strings.NewReader(content), opts...)

	return buf.String()
}

// Box writes the provided content in a simple box shape to given writer
func Box(out io.Writer, headline string, content io.Reader, opts ...BoxStyle) {
	var (
		beginning   = "╭"
		prefix      = "│"
		lastline    = "╵"
		linewritten = false
		outprint    = func(format string, a ...interface{}) {
			out.Write([]byte(fmt.Sprintf(format, a...)))
		}
	)

	// Process all provided box style options
	options := &boxOptions{}
	for _, f := range opts {
		f(options)
	}

	// Apply headline color if it is set
	if options.headlineColor != nil {
		for _, pointer := range []*string{&beginning, &headline, &prefix, &lastline} {
			*pointer = bunt.Style(*pointer,
				bunt.Foreground(*options.headlineColor),
			)
		}
	}

	// Apply headline styles if they are set
	for _, style := range options.headlineStyles {
		headline = bunt.Style(headline, style)
	}

	// Process each line of the content and apply styles if necessary
	scanner := bufio.NewScanner(content)
	for scanner.Scan() {
		text := scanner.Text()
		if options.contentColor != nil {
			text = bunt.Style(text, bunt.Foreground(*options.contentColor))
		}

		if !linewritten {
			// Write the headline string including the corner item
			outprint("%s %s\n", beginning, headline)
		}

		outprint("%s %s\n", prefix, text)
		linewritten = true
	}

	if linewritten {
		outprint("%s\n", lastline)
	}
}
