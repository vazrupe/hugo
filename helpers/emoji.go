// Copyright 2016 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helpers

import (
	"bytes"
	"index/suffixarray"
	"sort"
	"sync"

	"github.com/kyokomi/emoji"
)

var (
	emojiInit sync.Once

	emojis = make(map[string][]byte)

	emojiDelim     = []byte(":")
	emojiWordDelim = []byte(" ")
	emojiMaxSize   int
)

// Emoji returns the emojy given a key, e.g. ":smile:", nil if not found.
func Emoji(key string) []byte {
	emojiInit.Do(initEmoji)
	return emojis[key]
}

// Emojify "emojifies" the input source.
// Note that the input byte slice will be modified if needed.
// See http://www.emoji-cheat-sheet.com/
func Emojify(source []byte) []byte {
	emojiInit.Do(initEmoji)

	sa := suffixarray.New(source)

	var delimIdxs sort.IntSlice

	delimIdxs = sa.Lookup(emojiDelim, -1)
	delimIdxs.Sort()

	offset := 0
	for i := 0; i < len(delimIdxs)-1; i++ {

		s := offset + delimIdxs[i]
		e := offset + delimIdxs[i+1]
		if (e - s) > emojiMaxSize {
			continue
		}

		nexWordDelim := bytes.Index(source[s:e], emojiWordDelim)
		if nexWordDelim != -1 {
			continue
		}

		emojiKey := source[s : e+1]
		if emoji, ok := emojis[string(emojiKey)]; ok {
			source = append(source[:s], append(emoji, source[e+1:]...)...)

			i++
			offset += len(emoji) - len(emojiKey)
		}
	}

	return source
}

func initEmoji() {
	emojiMap := emoji.CodeMap()

	for k, v := range emojiMap {
		emojis[k] = []byte(v)

		if len(k) > emojiMaxSize {
			emojiMaxSize = len(k)
		}
	}

}
