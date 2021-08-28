// Copyright (c) 2017 Dave Pifke.
//
// Redistribution and use in source and binary forms, with or without
// modification, is permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

// Package wpasupplicant provides an interface for talking to the
// wpa_supplicant daemon.
//
// At the moment, this simply provides an interface for fetching wifi scan
// results.  More functionality is (probably) coming soon.
package wpasupplicant

// Cipher is one of the WPA_CIPHER constants from the wpa_supplicant source.
type Cipher int

const (
	CIPHER_NONE Cipher = 1 << iota
	WEP40
	WEP104
	TKIP
	CCMP
	AES_128_CMAC
	GCMP
	SMS4
	GCMP_256
	CCMP_256
	_
	BIP_GMAC_128
	BIP_GMAC_256
	BIP_CMAC_256
	GTK_NOT_USED
)

// KeyMgmt is one of the WPA_KEY_MGMT constants from the wpa_supplicant
// source.
type KeyMgmt int

const (
	IEEE8021X KeyMgmt = 1 << iota
	PSK
	KEY_MGMT_NONE
	IEEE8021X_NO_WPA
	WPA_NONE
	FT_IEEE8021X
	FT_PSK
	IEEE8021X_SHA256
	PSK_SHA256
	WPS
	SAE
	FT_SAE
	WAPI_PSK
	WAPI_CERT
	CCKM
	OSEN
	IEEE8021X_SUITE_B
	IEEE8021X_SUITE_B_192
)

type Algorithm int

type WPAEvent struct {
	Event     string
	Arguments map[string]string
	Line      string
}

// stdSocketPath is where to find the the AF_UNIX sockets for each interface.  It
// can be overridden for testing.
const stdSocketPath = "/run/wpa_supplicant"
