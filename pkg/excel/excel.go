package excel

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/xuri/excelize/v2"

	"github.com/koyashiro/md2xlsx/pkg/spec"
)

const CategoryCol = "A"
const SubCategoryCol = "B"
const SubSubCategoryCol = "C"
const ProceduresCol = "D"
const ConfirmationsCol = "E"
const ConfirmationPrefix = '・'

const CategoryWidth = 15.
const SubCategoryWidth = 15.
const SubSubCategoryWidth = 15.
const ProcedureWidth = 30.
const ConfirmationWidth = 30.

type Book struct {
	file *excelize.File
}

func NewBook() *Book {
	f := excelize.NewFile()
	f.DeleteSheet("sheet1")
	return &Book{file: f}
}

func OpenBook(filename string) (*Book, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	return &Book{file: f}, nil
}

func (b *Book) SaveAs(name string) error {
	return b.file.SaveAs(name)
}

func (b *Book) WriteTo(w io.Writer) (int64, error) {
	return b.file.WriteTo(w)
}

func (b *Book) WriteToBuffer() (*bytes.Buffer, error) {
	return b.file.WriteToBuffer()
}

func (b *Book) WriteSpec(spec *spec.Spec) error {
	var sheet string
	if spec.Name == "" {
		sheet = "no title"
	} else {
		sheet = spec.Name
	}
	b.file.NewSheet(sheet)
	b.file.DeleteSheet("sheet1")

	if err := setCelsWidth(b.file, sheet); err != nil {
		return err
	}

	sb := strings.Builder{}
	i := 2

	for _, c := range spec.Categories {
		if err := setCategory(b.file, sheet, i, c.Name); err != nil {
			return err
		}

		from := i
		for _, sc := range c.SubCategories {
			if err := setSubCategory(b.file, sheet, i, sc.Name); err != nil {
				return err
			}

			from := i
			for _, ssc := range sc.SubSubCategories {
				if err := setSubSubCategory(b.file, sheet, i, ssc.Name); err != nil {
					return err
				}
				if err := setConfirmations(b.file, sheet, i, ssc.Confirmations, &sb); err != nil {
					return err
				}
				if err := setProcedures(b.file, sheet, i, ssc.Procedures, &sb); err != nil {
					return err
				}
				i++
			}
			to := i - 1

			hcell := fmt.Sprintf("%s%d", SubCategoryCol, from)
			vcell := fmt.Sprintf("%s%d", SubCategoryCol, to)
			if err := b.file.MergeCell(sheet, hcell, vcell); err != nil {
				return err
			}
		}
		to := i - 1

		hcell := fmt.Sprintf("%s%d", CategoryCol, from)
		vcell := fmt.Sprintf("%s%d", CategoryCol, to)
		if err := b.file.MergeCell(sheet, hcell, vcell); err != nil {
			return err
		}
	}

	return nil
}

func setCelsWidth(f *excelize.File, sheet string) error {
	if err := f.SetColWidth(sheet, CategoryCol, CategoryCol, CategoryWidth); err != nil {
		return err
	}
	if err := f.SetColWidth(sheet, SubCategoryCol, SubCategoryCol, SubCategoryWidth); err != nil {
		return err
	}
	if err := f.SetColWidth(sheet, SubSubCategoryCol, SubSubCategoryCol, SubSubCategoryWidth); err != nil {
		return err
	}
	if err := f.SetColWidth(sheet, ProceduresCol, ProceduresCol, ProcedureWidth); err != nil {
		return err
	}
	if err := f.SetColWidth(sheet, ConfirmationsCol, ConfirmationsCol, ConfirmationWidth); err != nil {
		return err
	}
	return nil
}

func setCategory(f *excelize.File, sheet string, row int, name string) error {
	axis := fmt.Sprintf("%s%d", CategoryCol, row)
	return f.SetCellValue(sheet, axis, name)
}

func setSubCategory(f *excelize.File, sheet string, row int, name string) error {
	axis := fmt.Sprintf("%s%d", SubCategoryCol, row)
	return f.SetCellValue(sheet, axis, name)
}

func setSubSubCategory(f *excelize.File, sheet string, row int, name string) error {
	axis := fmt.Sprintf("%s%d", SubSubCategoryCol, row)
	return f.SetCellValue(sheet, axis, name)
}

func setProcedures(f *excelize.File, sheet string, row int, procedures []string, sb *strings.Builder) error {
	sb.Reset()
	axis := fmt.Sprintf("%s%d", ProceduresCol, row)
	for j, p := range procedures {
		if j != 0 {
			sb.WriteRune('\n')
		}
		sb.WriteString(fmt.Sprintf("%d. ", j+1))
		sb.WriteString(p)
	}
	return f.SetCellValue(sheet, axis, sb.String())
}

func setConfirmations(f *excelize.File, sheet string, row int, confirmations []string, sb *strings.Builder) error {
	sb.Reset()
	axis := fmt.Sprintf("%s%d", ConfirmationsCol, row)
	for j, p := range confirmations {
		if j != 0 {
			sb.WriteRune('\n')
		}
		sb.WriteRune(ConfirmationPrefix)
		sb.WriteString(p)
	}
	return f.SetCellValue(sheet, axis, sb.String())
}
