/**
 * Copyright 2023 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package abi

import (
    _ `unsafe`

    `github.com/bytedance/sonic/internal/rt`
)

const (
    _G_stackguard0 = 0x10
)

var (
    F_morestack_noctxt = uintptr(rt.FuncAddr(morestack_noctxt))
)

//go:linkname morestack_noctxt runtime.morestack_noctxt
func morestack_noctxt()

