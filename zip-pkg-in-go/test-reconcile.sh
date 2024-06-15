mkdir tmp-out
mkdir tmp-report
cp testdata/excel/pkg-test-xlsx-report*.xml   tmp-report

#go run main/main.go --command reconcile \
#                    --xls testdata/excel/pkg-test.xlsx \
#                    --config testdata/configs/test-pkg-in-400kb.properties \
#                    --report-dir tmp-report \
#                    --report-file-ends-with .xml \
#                    --out-dir tmp-out \
#                    --sheet-name Sheet1
#
#go run main/main.go --command reconcile \
#                    --xls testdata/excel/pkg-test.xlsx \
#                    --config testdata/configs/test-pkg-in-400kb.properties \
#                    --report-dir tmp-report \
#                    --report-file-ends-with .xml \
#                    --out-dir tmp-out \
#                    --sheet-name Sheet2

go run main/main.go --excel  testdata/excel/pkg-test.xlsx \
                    --sheet  Sheet1 \
                    --config testdata/configs/test-pkg-in-400kb.properties \
                    --out    tmp-out \
                    reconcile --report-dir tmp-report --report-file-ends-with .xml

go run main/main.go --excel  testdata/excel/pkg-test.xlsx \
                    --sheet  Sheet2 \
                    --config testdata/configs/test-pkg-in-400kb.properties \
                    --out    tmp-out \
                    reconcile --report-dir tmp-report --report-file-ends-with .xml
