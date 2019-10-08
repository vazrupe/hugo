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
	"sync"

	cedar "github.com/iohub/ahocorasick"
	"github.com/kyokomi/emoji"
)

var (
	emojiInit sync.Once

	emojis = make(map[string][]byte)

	matcher = cedar.NewMatcher()
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

	resp := matcher.Match(source)
	for resp.HasNext() {
		items := resp.NextMatchItem(source)
		for _, itr := range items {
			loc := itr.At - itr.KLen + 1

			source = append(source[:loc], append(itr.Value.([]byte), source[itr.At+1:]...)...)
		}
	}

	return source
}

func initEmoji() {
	emojiMap := emoji.CodeMap()

	for k, v := range emojiMap {
		emojis[k] = []byte(v)
		matcher.Insert([]byte(k), []byte(v))
	}

}
