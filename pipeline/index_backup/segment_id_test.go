/* Copyright © INFINI Ltd. All rights reserved.
 * web: https://infinilabs.com
 * mail: hello#infini.ltd */

package index_backup

import (
	"fmt"
	"testing"
)

func TestParseSegmentID(t *testing.T) {
	fileName:="_3g_Lucene85FieldsIndexfile_pointers_6x"
	fileName1:="_3e.fdt"

	//3g>3e

	segmentID1:=ParseSegmentID(fileName)
	fmt.Println(segmentID1)

	segmentID2:=ParseSegmentID(fileName1)
	fmt.Println(segmentID2)

	fmt.Println(CompareSegmentIDs(segmentID1,segmentID2))
	fmt.Println(CompareSegmentIDs(segmentID2,segmentID1))
	fmt.Println(CompareSegmentIDs(segmentID2,segmentID2))
	fmt.Println(CompareSegmentIDs("12","123"))

}