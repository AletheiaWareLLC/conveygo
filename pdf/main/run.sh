#!/bin/bash
#
# Copyright 2020 Aletheia Ware LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e
set -x

CORE_FONT_AFM_ZIP_URL=ftp://ftp.adobe.com/pub/adobe/devnet/font/pdfs/Core14_AFMs.zip
CORE_FONT_AFM_ZIP=Core14_AFMs.zip

if [ ! -f "$CORE_FONT_AFM_ZIP" ]; then
    curl -O "$CORE_FONT_AFM_ZIP_URL"
fi

go fmt $GOPATH/src/github.com/AletheiaWareLLC/{conveygo,conveygo/pdf,conveygo/pdf/graphics,conveygo/pdf/main}
go vet $GOPATH/src/github.com/AletheiaWareLLC/{conveygo,conveygo/pdf,conveygo/pdf/graphics,conveygo/pdf/main}
go test $GOPATH/src/github.com/AletheiaWareLLC/{conveygo,conveygo/pdf,conveygo/pdf/graphics,conveygo/pdf/main}
export ROOT_DIRECTORY=.
export ALIAS=test-convey-digest
export PASSWORD=password1234
go run github.com/AletheiaWareLLC/conveygo/pdf/main $@
