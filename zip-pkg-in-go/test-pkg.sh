mkdir tmp-src
mkdir tmp-out
cp testdata/pdfs/David-Passport.pdf      tmp-src/.
cp testdata/pdfs/Linda-DriverLicense.png tmp-src/.
cp testdata/pdfs/David-Passport.pdf      tmp-src/pp-0001.pdf
cp testdata/pdfs/Linda-DriverLicense.png tmp-src/dl-0001.pdf

#go run main/main.go --command package \
#                    --xls testdata/excel/pkg-test.xlsx \
#                    --file-dir tmp-src \
#                    --config testdata/configs/test-pkg-in-400kb.properties \
#                    --out-dir tmp-out \
#                    --unzip true \
#                    --sheet-name Sheet1
#
#go run main/main.go --command package \
#                    --xls testdata/excel/pkg-test.xlsx \
#                    --file-dir tmp-src \
#                    --config testdata/configs/test-pkg-in-3mb.properties \
#                    --out-dir tmp-out \
#                    --unzip true \
#                    --sheet-name Sheet1

#go run main/main.go --excel  testdata/excel/pkg-test.xlsx \
#                    --sheet  Sheet1 \
#                    --config testdata/configs/test-pkg-in-3mb.properties \
#                    --out    tmp-out \
#                    package \
#                    --pdf-dir tmp-src \
#                    --unzip-off

go run main/main.go --excel  testdata/excel/pkg-test.xlsx \
                    --sheet  Sheet1 \
                    --config testdata/configs/test-pkg-in-400kb.properties \
                    --out    tmp-out \
                    package \
                    --pdf-dir tmp-src
