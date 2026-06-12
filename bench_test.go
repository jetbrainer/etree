// Copyright 2015-2019 Brett Vickers.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package etree

import (
	"strconv"
	"strings"
	"testing"
)

var (
	benchmarkDocResult      *Document
	benchmarkElementResult  *Element
	benchmarkElementsResult []*Element
	benchmarkStringResult   string
	benchmarkBytesResult    []byte
	benchmarkCountResult    int
)

func benchmarkXML(sectionCount, bookCount int) string {
	var b strings.Builder
	b.Grow(sectionCount * bookCount * 180)
	b.WriteString(`<library>`)
	for section := 0; section < sectionCount; section++ {
		b.WriteString(`<section name="section`)
		b.WriteString(strconv.Itoa(section))
		b.WriteString(`">`)
		for book := 0; book < bookCount; book++ {
			category := "fiction"
			if book%3 == 0 {
				category = "web"
			}
			b.WriteString(`<book id="`)
			b.WriteString(strconv.Itoa(section))
			b.WriteString(`-`)
			b.WriteString(strconv.Itoa(book))
			b.WriteString(`" category="`)
			b.WriteString(category)
			b.WriteString(`"><title>Title `)
			b.WriteString(strconv.Itoa(book))
			b.WriteString(`</title><author>Author `)
			b.WriteString(strconv.Itoa(section))
			b.WriteString(`</author><year>`)
			b.WriteString(strconv.Itoa(2000 + book%25))
			b.WriteString(`</year><price currency="USD">`)
			b.WriteString(strconv.Itoa(10 + book%90))
			b.WriteString(`.00</price></book>`)
		}
		b.WriteString(`</section>`)
	}
	b.WriteString(`</library>`)
	return b.String()
}

func benchmarkDocument(sectionCount, bookCount int) *Document {
	doc := NewDocument()
	root := doc.CreateElement("library")
	for section := 0; section < sectionCount; section++ {
		sec := root.CreateElement("section")
		sec.CreateAttr("name", "section"+strconv.Itoa(section))
		for book := 0; book < bookCount; book++ {
			category := "fiction"
			if book%3 == 0 {
				category = "web"
			}
			item := sec.CreateElement("book")
			item.CreateAttr("id", strconv.Itoa(section)+"-"+strconv.Itoa(book))
			item.CreateAttr("category", category)
			item.CreateElement("title").SetText("Title " + strconv.Itoa(book))
			item.CreateElement("author").SetText("Author " + strconv.Itoa(section))
			item.CreateElement("year").SetText(strconv.Itoa(2000 + book%25))
			price := item.CreateElement("price")
			price.CreateAttr("currency", "USD")
			price.SetText(strconv.Itoa(10+book%90) + ".00")
		}
	}
	return doc
}

func benchmarkPrefixedDocument(sectionCount, bookCount int) *Document {
	doc := NewDocument()
	doc.WriteSettings = WriteSettings{
		CanonicalEndTags: true,
		CanonicalText:    true,
		CanonicalAttrVal: true,
	}
	root := doc.CreateElement("ds:library")
	root.CreateAttr("xmlns:ds", "urn:benchmark:ds")
	root.CreateAttr("xmlns:xsi", "http://www.w3.org/2001/XMLSchema-instance")
	for section := 0; section < sectionCount; section++ {
		sec := root.CreateElement("ds:section")
		sec.CreateAttr("ds:name", "section"+strconv.Itoa(section))
		sec.CreateAttr("xsi:type", "ds:Section")
		for book := 0; book < bookCount; book++ {
			item := sec.CreateElement("ds:book")
			item.CreateAttr("ds:id", strconv.Itoa(section)+"-"+strconv.Itoa(book))
			item.CreateAttr("ds:category", "web")
			item.CreateAttr("xsi:type", "ds:Book")
			item.CreateElement("ds:title").SetText("Title " + strconv.Itoa(book))
			item.CreateElement("ds:author").SetText("Author " + strconv.Itoa(section))
			price := item.CreateElement("ds:price")
			price.CreateAttr("ds:currency", "USD")
			price.SetText(strconv.Itoa(10+book%90) + ".00")
		}
	}
	return doc
}

func mustReadBenchmarkDocument(b *testing.B, xml string) *Document {
	b.Helper()
	doc := NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		b.Fatal(err)
	}
	return doc
}

func BenchmarkReadFromString(b *testing.B) {
	benchmarks := []struct {
		name string
		xml  string
	}{
		{name: "Small", xml: benchmarkXML(2, 5)},
		{name: "Medium", xml: benchmarkXML(20, 25)},
		{name: "Large", xml: benchmarkXML(80, 50)},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(bm.xml)))
			for i := 0; i < b.N; i++ {
				doc := NewDocument()
				if err := doc.ReadFromString(bm.xml); err != nil {
					b.Fatal(err)
				}
				benchmarkDocResult = doc
			}
		})
	}
}

func BenchmarkReadFromBytesMedium(b *testing.B) {
	input := []byte(benchmarkXML(20, 25))

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	for i := 0; i < b.N; i++ {
		doc := NewDocument()
		if err := doc.ReadFromBytes(input); err != nil {
			b.Fatal(err)
		}
		benchmarkDocResult = doc
	}
}

func BenchmarkWriteToStringMedium(b *testing.B) {
	input := benchmarkXML(20, 25)
	doc := mustReadBenchmarkDocument(b, input)

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, err := doc.WriteToString()
		if err != nil {
			b.Fatal(err)
		}
		benchmarkStringResult = s
	}
}

func BenchmarkWriteToStringPrefixed(b *testing.B) {
	doc := benchmarkPrefixedDocument(20, 25)
	input, err := doc.WriteToString()
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, err := doc.WriteToString()
		if err != nil {
			b.Fatal(err)
		}
		benchmarkStringResult = s
	}
}

func BenchmarkWriteToBytesMedium(b *testing.B) {
	input := benchmarkXML(20, 25)
	doc := mustReadBenchmarkDocument(b, input)

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := doc.WriteToBytes()
		if err != nil {
			b.Fatal(err)
		}
		benchmarkBytesResult = out
	}
}

func BenchmarkWriteToBytesPrefixed(b *testing.B) {
	doc := benchmarkPrefixedDocument(20, 25)
	input, err := doc.WriteToBytes()
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := doc.WriteToBytes()
		if err != nil {
			b.Fatal(err)
		}
		benchmarkBytesResult = out
	}
}

func BenchmarkIndentMedium(b *testing.B) {
	base := mustReadBenchmarkDocument(b, benchmarkXML(20, 25))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		doc := base.Copy()
		b.StartTimer()
		doc.Indent(2)
		benchmarkDocResult = doc
	}
}

func BenchmarkFindElementsMedium(b *testing.B) {
	doc := mustReadBenchmarkDocument(b, benchmarkXML(20, 25))
	path := ".//book[@category='web']/title"
	compiledPath := MustCompilePath(path)

	b.Run("StringPath", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkElementsResult = doc.FindElements(path)
		}
	})

	b.Run("CompiledPath", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkElementsResult = doc.FindElementsPath(compiledPath)
		}
	})

	b.Run("CompiledPathSeq", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			count := 0
			for range doc.FindElementsPathSeq(compiledPath) {
				count++
			}
			benchmarkCountResult = count
		}
	})
}

func BenchmarkFindElementMedium(b *testing.B) {
	doc := mustReadBenchmarkDocument(b, benchmarkXML(20, 25))
	path := ".//book[@category='web']/title"
	compiledPath := MustCompilePath(path)

	b.Run("StringPath", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkElementResult = doc.FindElement(path)
		}
	})

	b.Run("CompiledPath", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkElementResult = doc.FindElementPath(compiledPath)
		}
	})
}

func BenchmarkSelectElementsMedium(b *testing.B) {
	doc := mustReadBenchmarkDocument(b, benchmarkXML(20, 25))
	root := doc.Root()

	b.Run("Slice", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			benchmarkElementsResult = root.SelectElements("section")
		}
	})

	b.Run("Seq", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			count := 0
			for range root.SelectElementsSeq("section") {
				count++
			}
			benchmarkCountResult = count
		}
	})
}

func BenchmarkCopyMedium(b *testing.B) {
	doc := mustReadBenchmarkDocument(b, benchmarkXML(20, 25))

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkDocResult = doc.Copy()
	}
}

func BenchmarkCreateDocumentMedium(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkDocResult = benchmarkDocument(20, 25)
	}
}
