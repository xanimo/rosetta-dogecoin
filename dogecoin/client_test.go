// Copyright 2020 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dogecoin

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

const (
	url = "/"
)

func forceMarshalMap(t *testing.T, i interface{}) map[string]interface{} {
	m, err := types.MarshalMap(i)
	if err != nil {
		t.Fatalf("could not marshal map %s", types.PrintStruct(i))
	}

	return m
}

var (
	blockIdentifier1000 = &types.BlockIdentifier{
		Hash:  "6da6ef0eb0e2e1e150ed44c43f596b9552b221a4dda207535eca9c0ddb7a10d0",
		Index: 1000,
	}

	block1000 = &Block{
		Hash:              "6da6ef0eb0e2e1e150ed44c43f596b9552b221a4dda207535eca9c0ddb7a10d0",
		Height:            1000,
		PreviousBlockHash: "172024fe6ffe31b206e159319ed391f6a532a69007c047e46d8bea06f89aae60",
		Time:              1386481098,
		Size:              190,
		Weight:            760,
		Version:           1,
		MerkleRoot:        "9480ac41aac2674ac498849b9ab95661ee73fb372140e62c6e7a6fa29f5a09d1",
		MedianTime:        1386480992,
		Nonce:             2308638208,
		Bits:              "1d03d07b",
		Difficulty:        0.2621620216098152,
		Txs: []*Transaction{
			{
				Hex:      "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0e04ca05a4520131062f503253482fffffffff01005637cfe7100000232103236a82a90eda514a3373f1a13b349cbff983d00879023ce15ef5ff8c757aa7dbac00000000", // nolint
				Hash:     "9480ac41aac2674ac498849b9ab95661ee73fb372140e62c6e7a6fa29f5a09d1",
				Size:     109,
				Vsize:    109,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						Coinbase: "04ca05a4520131062f503253482f",
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 185878.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:  "03236a82a90eda514a3373f1a13b349cbff983d00879023ce15ef5ff8c757aa7db OP_CHECKSIG", // nolint
							Hex:  "2103236a82a90eda514a3373f1a13b349cbff983d00879023ce15ef5ff8c757aa7dbac",         // nolint
							Type: "pubkey",
							Addresses: []string{
								"DG3FAujuBTvozAU423arUKYnhQQa6XzWZ9",
							},
						},
					},
				},
			},
		},
	}

	blockIdentifier100000 = &types.BlockIdentifier{
		Hash:  "000000000003ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506",
		Index: 100000,
	}

	block100000 = &Block{
		Hash:              "13ab3b961fcc500c03f51279385c42e9f055d48a37dfa72d0073c0d3f595036b",
		Height:            100000,
		PreviousBlockHash: "12aca0938fe1fb786c9e0e4375900e8333123de75e240abd3337d1b411d14ebe",
		Time:              1392346781,
		Size:              13372,
		Weight:            53488,
		Version:           1,
		MerkleRoot:        "31757c266102d1bee62ef2ff8438663107d64bdd5d9d9173421ec25fb2a814de",
		MedianTime:        1392346434,
		Nonce:             2216773632,
		Bits:              "1b267eeb",
		Difficulty:        1702.39468793143,
		Txs: []*Transaction{
			{
				Hex:      "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff2303a08601062f503253482f049186fd5208080030c207000000092f7374726174756d2f0000000001a8f80742a91f00001976a9146209ed8017ee7c7efcfccb5c971ba58a951749af88ac00000000", // nolint
				Hash:     "c2e410a0c9ff9ee2808a5efc27885e092baf337290dca77faf1b13da5a946d98",
				Size:     120,
				Vsize:    120,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						Coinbase: "03a08601062f503253482f049186fd5208080030c207000000092f7374726174756d2f",
						Sequence: 0,
					},
				},
				Outputs: []*Output{
					{
						Value: 348118.17752744,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 6209ed8017ee7c7efcfccb5c971ba58a951749af OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9146209ed8017ee7c7efcfccb5c971ba58a951749af88ac",
							RequiredSigs: 1,
							Type:         "scripthash",
							Addresses: []string{
								"DE5UerCH9YNvVDnj3Firj4sKrAex2krgnt",
							},
						},
					},
				},
			},
			{
				Hex:      "010000000116b362b652c601435f94a6ecfc6c71a1fd240d58f5e4443cc6f3f31f7eae500b010000006b483045022100ca9176c3eccd6ab443b1259698c35a40f8274e5b87222e2a69a1d1937398c59f02205390042fa802df9a26923e6d972ed837cb7f297e307cb84a3cdee61c8d6f3a08012102f8ae61694000cff50ae14e80994c34fa6bb672d8503e904adc0f43dd7ec14f04ffffffff0200a0724e180900001976a914101f0445d2cee10c1f820dac0fdab961c35c594088acf2e893c4a45801001976a9146026526eeaff25b0dbf335e95241ff7ff990a59888ac00000000", // nolint
				Hash:     "19220173c151925d493a25dbe67798fa11e6e4db01b927b1c04cddead85d3d12",
				Size:     226,
				Vsize:    226,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "0b50ae7e1ff3f3c63c44e4f5580d24fda1716cfceca6945f4301c652b662b316",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3045022100ca9176c3eccd6ab443b1259698c35a40f8274e5b87222e2a69a1d1937398c59f02205390042fa802df9a26923e6d972ed837cb7f297e307cb84a3cdee61c8d6f3a08[ALL] 02f8ae61694000cff50ae14e80994c34fa6bb672d8503e904adc0f43dd7ec14f04", // nolint
							Hex: "483045022100ca9176c3eccd6ab443b1259698c35a40f8274e5b87222e2a69a1d1937398c59f02205390042fa802df9a26923e6d972ed837cb7f297e307cb84a3cdee61c8d6f3a08012102f8ae61694000cff50ae14e80994c34fa6bb672d8503e904adc0f43dd7ec14f04", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 100000.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 101f0445d2cee10c1f820dac0fdab961c35c5940 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914101f0445d2cee10c1f820dac0fdab961c35c594088ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D6cLWQTCnvCxUZS5G99dsDmKc3vGuo26w6",
							},
						},
					},
					{
						Value: 3789396.72619250,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 6026526eeaff25b0dbf335e95241ff7ff990a598 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9146026526eeaff25b0dbf335e95241ff7ff990a59888ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DDuVKLJyHqpPBKZGVt4wREo3c6tR4GaafX",
							},
						},
					},
				},
			},

			{
				Hex:      "0100000001d405cda6ec490b6189a09f22d3fa68795b91563aca04a91aa1964b70b2a0a674000000006a47304402200ef30fb2d38f4f28bcee60e9e5f9b0d45a09fc908ddf459eec231d128a2f14540220683d7eade3a0cce01ac986eed294cf8644cf0689ae67823853dca88591d98cdd0121039e8e4887fc41bbf3c2384548391b76e3e43ab8880416947f640b11fc665cbef5ffffffff02ce25bc125bd406001976a9140ca60499e9ccf01717dfd7a9c40f7bc8226a448f88ac0ff005b3520000001976a9142ce991b0fa2bb82b085382685c90b8f7f19513db88ac00000000", // nolint
				Hash:     "7da232b380f7a0c05bcb5e526e17be7d156c0e6686ce0e72df4cdd1e2afb6458",
				Size:     225,
				Vsize:    225,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "74a6a0b2704b96a11aa904ca3a56915b7968fad3229fa089610b49eca6cd05d4",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "304402200ef30fb2d38f4f28bcee60e9e5f9b0d45a09fc908ddf459eec231d128a2f14540220683d7eade3a0cce01ac986eed294cf8644cf0689ae67823853dca88591d98cdd[ALL] 039e8e4887fc41bbf3c2384548391b76e3e43ab8880416947f640b11fc665cbef5", // nolint
							Hex: "47304402200ef30fb2d38f4f28bcee60e9e5f9b0d45a09fc908ddf459eec231d128a2f14540220683d7eade3a0cce01ac986eed294cf8644cf0689ae67823853dca88591d98cdd0121039e8e4887fc41bbf3c2384548391b76e3e43ab8880416947f640b11fc665cbef5", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 19223374.81696718,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 0ca60499e9ccf01717dfd7a9c40f7bc8226a448f OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9140ca60499e9ccf01717dfd7a9c40f7bc8226a448f88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D6HyXh4Dv1zs2QQh9fjaDLbdQkRY4rzTB3",
							},
						},
					},
					{
						Value: 3551.90829071,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 2ce991b0fa2bb82b085382685c90b8f7f19513db OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9142ce991b0fa2bb82b085382685c90b8f7f19513db88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D9Ea5N8BNa3NcxdkwueYUjj83LWzpK2jvc",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000013f263cc3b5077e90920396055357b2d44535b222cff61579211cbaac51f388c8000000006a47304402206b12cd8abfea1a24cf90136b46f57d39cf0c826a5ec070296396ea2bc2c13152022054bf823e9dc35cb09867736189f31df1f51ce4cd6fcff06be51db8ba5e5ffd8501210287614929460d7acb7aef3b11194fff34e20d6253842d1f73e535a4fe7b5f70dfffffffff0276272c6fa30000001976a91427691601ceb91d590c42d52e6e90868ab2e5f23588ac00e87648170000001976a914ab215e1613facb3a426a4f553e6c7b20b58bf86088ac00000000", // nolint
				Hash:     "e79451d5d86a6919b01bfae9e00bcbccdf7b8b38099dcb38336e18e6b879684c",
				Size:     225,
				Vsize:    225,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "c888f351acba1c217915f6cf22b23545d4b2575305960392907e07b5c33c263f",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "304402206b12cd8abfea1a24cf90136b46f57d39cf0c826a5ec070296396ea2bc2c13152022054bf823e9dc35cb09867736189f31df1f51ce4cd6fcff06be51db8ba5e5ffd85[ALL] 0287614929460d7acb7aef3b11194fff34e20d6253842d1f73e535a4fe7b5f70df", // nolint
							Hex: "47304402206b12cd8abfea1a24cf90136b46f57d39cf0c826a5ec070296396ea2bc2c13152022054bf823e9dc35cb09867736189f31df1f51ce4cd6fcff06be51db8ba5e5ffd8501210287614929460d7acb7aef3b11194fff34e20d6253842d1f73e535a4fe7b5f70df", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 7019.44833910,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 27691601ceb91d590c42d52e6e90868ab2e5f235 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a91427691601ceb91d590c42d52e6e90868ab2e5f23588ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D8jUnEuuuRyJJxrenFQwz6pH8thvY7ytGD",
							},
						},
					},
					{
						Value: 1000.00000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 ab215e1613facb3a426a4f553e6c7b20b58bf860 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914ab215e1613facb3a426a4f553e6c7b20b58bf86088ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DLjx4UfAm8Dm3JKNQ9WuQcgb7THomXDupw",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000020045a1d9128ebf682cccae605df74316eb1b2644f656143866cad5f1ceebd99a010000006b4830450221008e71a618a5c424667664b52401fd42c2ce38f41b29a9250d72fa3781b1b4585c02200cde3f5ad65e61d5a8e792a2cd6f742ee0a8d08d6a194293b92940d85285cb080121030e2247d028860a11b4ca970f81dfac88a0dbecb44d4f11100045ce365b6155e1ffffffff33fafe002f76a2ddcb1862762eb6bffd75ad8d549277bdb047af211153187026010000006b483045022100f39be19e4786c4da61e58b46812b56e8b53ccb533a7dddc7695ae82ebb1e2bef02203fabde51827acb8fc918df241c02856bc3be01491c2aea9f07f4ea19aa78e2a40121030e2247d028860a11b4ca970f81dfac88a0dbecb44d4f11100045ce365b6155e1ffffffff01c1d2c57fec0100001976a914cf4f2986e58650a5acece7c81a56a374c9c5170188ac00000000", // nolint
				Hash:     "f50856b0a7cb9893c6a1e9d2d17199c369fc5752d68ed7d5651bda9e19b8ba9b",
				Size:     340,
				Vsize:    340,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "9ad9ebcef1d5ca66381456f644261beb1643f75d60aecc2c68bf8e12d9a14500",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "30450221008e71a618a5c424667664b52401fd42c2ce38f41b29a9250d72fa3781b1b4585c02200cde3f5ad65e61d5a8e792a2cd6f742ee0a8d08d6a194293b92940d85285cb08[ALL] 030e2247d028860a11b4ca970f81dfac88a0dbecb44d4f11100045ce365b6155e1", // nolint
							Hex: "4830450221008e71a618a5c424667664b52401fd42c2ce38f41b29a9250d72fa3781b1b4585c02200cde3f5ad65e61d5a8e792a2cd6f742ee0a8d08d6a194293b92940d85285cb080121030e2247d028860a11b4ca970f81dfac88a0dbecb44d4f11100045ce365b6155e1", // nolint
						},
						Sequence: 4294967295,
					},
					{
						TxHash: "267018531121af47b0bd7792548dad75fdbfb62e766218cbdda2762f00fefa33",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "3045022100f39be19e4786c4da61e58b46812b56e8b53ccb533a7dddc7695ae82ebb1e2bef02203fabde51827acb8fc918df241c02856bc3be01491c2aea9f07f4ea19aa78e2a4[ALL] 030e2247d028860a11b4ca970f81dfac88a0dbecb44d4f11100045ce365b6155e1", // nolint
							Hex: "483045022100f39be19e4786c4da61e58b46812b56e8b53ccb533a7dddc7695ae82ebb1e2bef02203fabde51827acb8fc918df241c02856bc3be01491c2aea9f07f4ea19aa78e2a40121030e2247d028860a11b4ca970f81dfac88a0dbecb44d4f11100045ce365b6155e1", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 21152.67580609,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 cf4f2986e58650a5acece7c81a56a374c9c51701 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914cf4f2986e58650a5acece7c81a56a374c9c5170188ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DQ3FFFiAkjPhthedzFW6eTE1kg6fdJrC1y",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000016b3a866916195b44b1ef1eeca6bd42444d30bd801653eae356827286d7170732000000006b48304502210095052e425fc05be614f6312bdcd8acb7ada011a1d3575d5dbcdd203c4ee09e6f022001f66b82f8c59c13a60f6aa84f5a991dd54cc47a89d189326e813665e7e7cb6e012103db50768756c7daaa9096e537afe542415da9649668023c26b9013ccb04a72db2ffffffff02711af552470e00001976a914301ab2c3a86972dee0bc2250be7e1e4f3f55886388aca3dede6fe20000001976a9142f119fccae2f2ea71689280aeec457b10fdd6e4488ac00000000", // nolint
				Hash:     "1458cf1ea545950631048c77c38f6a4c12a96b51164864ac1e097cfc18c28241",
				Size:     226,
				Vsize:    226,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "320717d786728256e3ea531680bd304d4442bda6ec1eefb1445b191669863a6b",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "304502210095052e425fc05be614f6312bdcd8acb7ada011a1d3575d5dbcdd203c4ee09e6f022001f66b82f8c59c13a60f6aa84f5a991dd54cc47a89d189326e813665e7e7cb6e[ALL] 03db50768756c7daaa9096e537afe542415da9649668023c26b9013ccb04a72db2", // nolint
							Hex: "48304502210095052e425fc05be614f6312bdcd8acb7ada011a1d3575d5dbcdd203c4ee09e6f022001f66b82f8c59c13a60f6aa84f5a991dd54cc47a89d189326e813665e7e7cb6e012103db50768756c7daaa9096e537afe542415da9649668023c26b9013ccb04a72db2", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 156994.97261681,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 301ab2c3a86972dee0bc2250be7e1e4f3f558863 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914301ab2c3a86972dee0bc2250be7e1e4f3f55886388ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D9XSxVkhaT5957H3vLw6sCSwyfXXWbTyjq",
							},
						},
					},
					{
						Value: 9725.39485859,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 2f119fccae2f2ea71689280aeec457b10fdd6e44 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9142f119fccae2f2ea71689280aeec457b10fdd6e4488ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D9RyQpthPE45xmA4pTGGS9APNWcWC4yps9",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000004c61f35d758b2dd4c8aa93161a7e6597b078545942219c9e1eb9d2560206a385b010000006b483045022100d758b832dd4a379ca2ce15f23ea42263b47608e483316ee7bad9153ac7a9e2f202205ab6e79d3147042945e6e1710fd0c571b0a7b9c606a257f34bfbec8d1d17786f012103557da6fd76f671f830b68f9d52976e720dde17ba43c4d09575c017f6bddb3971ffffffff3b124474a1f49e26b62314d1ae2c72222b6a2022191b6733989382a844d386f6010000006a473044022036de6e36f40a8db90e565fc398b463742d6e6751b6c5eb758563138f4a063e3a022078234eda7e829eb9458c7ab35a210fa47f1355d48c61ff5cfb31be7e42ed33e001210303911d4a1c0ad1530ad241ca2caab9038b21ed86055ac940ea58e5634c385501ffffffff5b953461686be8fd321b14f90ec8f3bc5e1b83063e2bf2894345d4a90aad0022010000006b483045022100d98d5f9dbec259b2fce4e9f8201b07f27f4025853beb331be91e1704e304b2ba02205808df049a59145f9526cc37cbed772468e58c3d7bc50a37c752456681c7bfd501210303911d4a1c0ad1530ad241ca2caab9038b21ed86055ac940ea58e5634c385501ffffffff0046dde73f58ddeb824e8997f39a1b417bb45d5f69a16e57d65830cc0d3bd8a7000000006c4930460221008a2dfb6cdaeda83eaddc3e4e79b5f5a4b63b813470f7cfc1d4e380f0d45a5dd9022100e05932debf59fceebf7d0fb3b55e50eaac1746b2b3ed27a3a446d92579937143012102e3c99ec61db654fcb1829fdfb707f475b501829a945fe46d994db2cc24fd9b10ffffffff010012ca97310000001976a914c737ab4785a1e88d8b6152e7207fc3372ff61ee088ac00000000", // nolint
				Hash:     "37ea455a9f80e9d40e2a4009ff27c600ac3bfbfd22a5143ebb2ce41b11a15ec0",
				Size:     636,
				Vsize:    636,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "5b386a2060259debe1c91922944585077b59e6a76131a98a4cddb258d7351fc6",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "3045022100d758b832dd4a379ca2ce15f23ea42263b47608e483316ee7bad9153ac7a9e2f202205ab6e79d3147042945e6e1710fd0c571b0a7b9c606a257f34bfbec8d1d17786f[ALL] 03557da6fd76f671f830b68f9d52976e720dde17ba43c4d09575c017f6bddb3971", // nolint
							Hex: "483045022100d758b832dd4a379ca2ce15f23ea42263b47608e483316ee7bad9153ac7a9e2f202205ab6e79d3147042945e6e1710fd0c571b0a7b9c606a257f34bfbec8d1d17786f012103557da6fd76f671f830b68f9d52976e720dde17ba43c4d09575c017f6bddb3971", // nolint
						},
						Sequence: 4294967295,
					},
					{
						TxHash: "f686d344a882939833671b1922206a2b22722caed11423b6269ef4a17444123b",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "3044022036de6e36f40a8db90e565fc398b463742d6e6751b6c5eb758563138f4a063e3a022078234eda7e829eb9458c7ab35a210fa47f1355d48c61ff5cfb31be7e42ed33e0[ALL] 0303911d4a1c0ad1530ad241ca2caab9038b21ed86055ac940ea58e5634c385501", // nolint
							Hex: "473044022036de6e36f40a8db90e565fc398b463742d6e6751b6c5eb758563138f4a063e3a022078234eda7e829eb9458c7ab35a210fa47f1355d48c61ff5cfb31be7e42ed33e001210303911d4a1c0ad1530ad241ca2caab9038b21ed86055ac940ea58e5634c385501", // nolint
						},
						Sequence: 4294967295,
					},
					{
						TxHash: "2200ad0aa9d4454389f22b3e06831b5ebcf3c80ef9141b32fde86b686134955b",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "3045022100d98d5f9dbec259b2fce4e9f8201b07f27f4025853beb331be91e1704e304b2ba02205808df049a59145f9526cc37cbed772468e58c3d7bc50a37c752456681c7bfd5[ALL] 0303911d4a1c0ad1530ad241ca2caab9038b21ed86055ac940ea58e5634c385501", // nolint
							Hex: "483045022100d98d5f9dbec259b2fce4e9f8201b07f27f4025853beb331be91e1704e304b2ba02205808df049a59145f9526cc37cbed772468e58c3d7bc50a37c752456681c7bfd501210303911d4a1c0ad1530ad241ca2caab9038b21ed86055ac940ea58e5634c385501", // nolint
						},
						Sequence: 4294967295,
					},
					{
						TxHash: "a7d83b0dcc3058d6576ea1695f5db47b411b9af397894e82ebdd583fe7dd4600",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "30460221008a2dfb6cdaeda83eaddc3e4e79b5f5a4b63b813470f7cfc1d4e380f0d45a5dd9022100e05932debf59fceebf7d0fb3b55e50eaac1746b2b3ed27a3a446d92579937143[ALL] 02e3c99ec61db654fcb1829fdfb707f475b501829a945fe46d994db2cc24fd9b10", // nolint
							Hex: "4930460221008a2dfb6cdaeda83eaddc3e4e79b5f5a4b63b813470f7cfc1d4e380f0d45a5dd9022100e05932debf59fceebf7d0fb3b55e50eaac1746b2b3ed27a3a446d92579937143012102e3c99ec61db654fcb1829fdfb707f475b501829a945fe46d994db2cc24fd9b10", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 2130.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 c737ab4785a1e88d8b6152e7207fc3372ff61ee0 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914c737ab4785a1e88d8b6152e7207fc3372ff61ee088ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DPJThPPSYWA5b64YSXBbRE5eJpqRVQbjte",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000012a3acddcefc9eeec8c3e3ad2842d8d1bd45285a89f30602b115678fba2ffdf0e000000006a4730440220263463850fda3f7e45c1f0926c0d5b9e16bb4fd5535eb8c23b77363f673973d702202b93f4623898f4e066e05077c6bc13e9122c427de55a329f54521927d9e6ee6a0121034bfaa8a15fe5d3384352580f1a84f7f49924115af4a7c8b02da899047c2314c4ffffffff02b61ffe532f0c00001976a91423c97abb9080b029155d4634f988783a00510ddc88ac08e2c7e0060000001976a9149b65a61c3413185d315ff12f326a9c3ac8a4debd88ac00000000", // nolint
				Hash:     "c7d255d3b6851c52a6d367ca237c7f02a6c85bfcf3af22e206d053a363888225",
				Size:     225,
				Vsize:    225,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "0edfffa2fb7856112b60309fa88552d41b8d2d84d23a3e8ceceec9efdccd3a2a",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "30440220263463850fda3f7e45c1f0926c0d5b9e16bb4fd5535eb8c23b77363f673973d702202b93f4623898f4e066e05077c6bc13e9122c427de55a329f54521927d9e6ee6a[ALL] 034bfaa8a15fe5d3384352580f1a84f7f49924115af4a7c8b02da899047c2314c4", // nolint
							Hex: "4730440220263463850fda3f7e45c1f0926c0d5b9e16bb4fd5535eb8c23b77363f673973d702202b93f4623898f4e066e05077c6bc13e9122c427de55a329f54521927d9e6ee6a0121034bfaa8a15fe5d3384352580f1a84f7f49924115af4a7c8b02da899047c2314c4", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 133974.12159414,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 23c97abb9080b029155d4634f988783a00510ddc OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a91423c97abb9080b029155d4634f988783a00510ddc88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D8QKZ1SJbKX6J2e24Jt2dA4z5fivEZUWuB",
							},
						},
					},
					{
						Value: 295.40999688,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 9b65a61c3413185d315ff12f326a9c3ac8a4debd OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9149b65a61c3413185d315ff12f326a9c3ac8a4debd88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DKJm3gMTBZ6ciDeQzmjiW6efQxmJvaVjHi",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000012a3acddcefc9eeec8c3e3ad2842d8d1bd45285a89f30602b115678fba2ffdf0e000000006a4730440220263463850fda3f7e45c1f0926c0d5b9e16bb4fd5535eb8c23b77363f673973d702202b93f4623898f4e066e05077c6bc13e9122c427de55a329f54521927d9e6ee6a0121034bfaa8a15fe5d3384352580f1a84f7f49924115af4a7c8b02da899047c2314c4ffffffff02b61ffe532f0c00001976a91423c97abb9080b029155d4634f988783a00510ddc88ac08e2c7e0060000001976a9149b65a61c3413185d315ff12f326a9c3ac8a4debd88ac00000000", // nolint
				Hash:     "74017e45e9eb84742d5ec34a9d9f670f90e4caceaaa4a00948233df0b6f039b9",
				Size:     227,
				Vsize:    227,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "5348802042c21c2b843c5f919b2286661c30d0ac32fd2b63d53a1b643c636a77",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3046022100e82121319e327d1070bd1b8afeeb060e88cbe10bd468d228dc77bb506f8b661d022100d00a705efd8265ef18113acf35ac9fbb46fe113902828e1d72eb4bcb9bdf7b6b[ALL] 036eb7d2dabc6560354251f252a51475c5ed559830681ba0838d89d494282d0c79", // nolint
							Hex: "493046022100e82121319e327d1070bd1b8afeeb060e88cbe10bd468d228dc77bb506f8b661d022100d00a705efd8265ef18113acf35ac9fbb46fe113902828e1d72eb4bcb9bdf7b6b0121036eb7d2dabc6560354251f252a51475c5ed559830681ba0838d89d494282d0c79", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 940164.49521597,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 42f185de5643d0e4653d9f50b54eec6d016769fe OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a91442f185de5643d0e4653d9f50b54eec6d016769fe88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D8QKZ1SJbKX6J2e24Jt2dA4z5fivEZUWuB",
							},
						},
					},
					{
						Value: 295.40999688,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 9b65a61c3413185d315ff12f326a9c3ac8a4debd OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914301ab2c3a86972dee0bc2250be7e1e4f3f55886388ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DKJm3gMTBZ6ciDeQzmjiW6efQxmJvaVjHi",
							},
						},
					},
				},
			},
			{
				Hex:      "010000000186befd244b4ac06a0c93f3c84b31ffe8364aa270c359f49c24a27e38aa37542d000000006a4730440220401773baf69df3b93fb0b51c82b2c05e39ea665e3302d060d1005368a06143c102206535d87af7aed9886d68682b082a408e52f41ea1a92d4aca90cec83610af8a9a012103aa2c2048bf90b9a2de755412c6b0c0f66f21dc76cc14bcc4f64cacc3e515c4bdffffffff02ae70d1b54f4900001976a91450008196da21b76532b492ab97dc991456ad179a88ac521d8077260000001976a9146d0bc9d3c908d759b0e0a59a5c26d69e37256c1f88ac00000000", // nolint
				Hash:     "79c7dbc959daad2099efdd2465813c045f945f79751e332e7b7b2c7972be5203",
				Size:     225,
				Vsize:    225,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "2d5437aa387ea2249cf459c370a24a36e8ff314bc8f3930c6ac04a4b24fdbe86",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "30440220401773baf69df3b93fb0b51c82b2c05e39ea665e3302d060d1005368a06143c102206535d87af7aed9886d68682b082a408e52f41ea1a92d4aca90cec83610af8a9a[ALL] 03aa2c2048bf90b9a2de755412c6b0c0f66f21dc76cc14bcc4f64cacc3e515c4bd", // nolint
							Hex: "4730440220401773baf69df3b93fb0b51c82b2c05e39ea665e3302d060d1005368a06143c102206535d87af7aed9886d68682b082a408e52f41ea1a92d4aca90cec83610af8a9a012103aa2c2048bf90b9a2de755412c6b0c0f66f21dc76cc14bcc4f64cacc3e515c4bd", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 806067.01645998,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 50008196da21b76532b492ab97dc991456ad179a OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a91450008196da21b76532b492ab97dc991456ad179a88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DCS7Cn1VoVBkoz3ESZrrN6AKFQN8neGjSb",
							},
						},
					},
					{
						Value: 1652.13642066,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 6d0bc9d3c908d759b0e0a59a5c26d69e37256c1f OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9146d0bc9d3c908d759b0e0a59a5c26d69e37256c1f88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DF5gKEqwG3QRAw6w7y6AA1bw3yaRqnRYUP",
							},
						},
					},
				},
			},
			{
				Hex:      "010000000137ed67ebea136256834f48ffed22e469cf5e5659d52227ff28a5213482defdb7000000006c493046022100907a2c31acd9ed8cab0f2b8c94249d1a5bfd06c01024cefddfd2225f967d32ab022100f9cb99ddb4edc5eebe70d5f5588767a978e21a7fb8c86c037ed6c1ef83e63de20121038cb41f889e7658adee33cb1d8ec5f9586424ed2a4b92e56373e756b1aad300aeffffffff02a4ed016fa20b00001976a9148841590909747c0f97af158f22fadacb1652522088ac409d1a9c0a0000001976a914efb6158f75743c611858fdfd0f4aaec6cc6196bc88ac00000000", // nolint
				Hash:     "2d607f80be822c73a4769cf58481896984e4f6c9cd596ba75ec47212837ca2c4",
				Size:     227,
				Vsize:    227,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "b7fdde823421a528ff2722d559565ecf69e422edff484f83566213eaeb67ed37",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3046022100907a2c31acd9ed8cab0f2b8c94249d1a5bfd06c01024cefddfd2225f967d32ab022100f9cb99ddb4edc5eebe70d5f5588767a978e21a7fb8c86c037ed6c1ef83e63de2[ALL] 038cb41f889e7658adee33cb1d8ec5f9586424ed2a4b92e56373e756b1aad300ae", // nolint
							Hex: "493046022100907a2c31acd9ed8cab0f2b8c94249d1a5bfd06c01024cefddfd2225f967d32ab022100f9cb99ddb4edc5eebe70d5f5588767a978e21a7fb8c86c037ed6c1ef83e63de20121038cb41f889e7658adee33cb1d8ec5f9586424ed2a4b92e56373e756b1aad300ae", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 100000.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 8841590909747c0f97af158f22fadacb16525220 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9148841590909747c0f97af158f22fadacb1652522088ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DHZYinsaM9nW5piCMN639ELRKbZomThPnZ",
							},
						},
					},
					{
						Value: 455.68662848,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 efb6158f75743c611858fdfd0f4aaec6cc6196bc OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914efb6158f75743c611858fdfd0f4aaec6cc6196bc88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DSzaAYEYyy9ngjoJ294r7jzFM3xhD6bKHK",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000015ac6959f5f8e0d14df1c45da8567d0a2dc4cf682431359cab2a192c81a854e3f000000006b48304502204f4132692261b237ed877b76a722a7e58d48e9157f013fbba6f6adc6944e83ce022100a0a64c156aabd409fea754d225a8c16e6a2a2920b1b6b3f01162734e160dcd530121022b9fafc77f470ed11756a82e4d27dc8229196779525b73aa170d0adde7ef050fffffffff02008dff52da0900001976a9142f6c6ee207403395836335e4460e96e7f583701c88ac00b080f6450100001976a9143eda58dd12a5edca008ce105fe3847f96016f17188ac00000000", // nolint
				Hash:     "ba9b4c98254e8215fc66535c4034a27cf731eaec5e9255720ca563969e36a79a",
				Size:     226,
				Vsize:    226,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "3f4e851ac892a1b2ca59134382f64cdca2d06785da451cdf140d8e5f9f95c65a",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "304502204f4132692261b237ed877b76a722a7e58d48e9157f013fbba6f6adc6944e83ce022100a0a64c156aabd409fea754d225a8c16e6a2a2920b1b6b3f01162734e160dcd53[ALL] 022b9fafc77f470ed11756a82e4d27dc8229196779525b73aa170d0adde7ef050f", // nolint
							Hex: "48304502204f4132692261b237ed877b76a722a7e58d48e9157f013fbba6f6adc6944e83ce022100a0a64c156aabd409fea754d225a8c16e6a2a2920b1b6b3f01162734e160dcd530121022b9fafc77f470ed11756a82e4d27dc8229196779525b73aa170d0adde7ef050f", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 108333.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 2f6c6ee207403395836335e4460e96e7f583701c OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9142f6c6ee207403395836335e4460e96e7f583701c88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D9TrCLiwDR1VKhNaa6YYXMGbdzmTMcik3Y",
							},
						},
					},
					{
						Value: 14000.00000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 3eda58dd12a5edca008ce105fe3847f96016f171 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9143eda58dd12a5edca008ce105fe3847f96016f17188ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DAsS1BANzXExDRp1wke4rKp9UFga5T9wGs",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000010daae87a56b318b250d9ae5c59573cdc555c5d1e73a2b6840b04755c1c2aba22010000006b4830450220720793be4b7812a7b4b7ae6dd75cf37da09ddd4630debc09c34bf95782d16bdc022100f5f609629a559bcd96dcf82cb1657aa6374af78d3114c2f57d9f6fb204c97718012103c584550e0b922323d39a357da5ee977e14a117fda1812ece91f33b06072cf49dffffffff0200943577000000001976a9142deeef21adf2cf1d0fc8cf31efcf11af2321595d88acba4c55c9000000001976a9148a10bf7df695647d408a07d765699d41280663d888ac00000000", // nolint
				Hash:     "904b55a5ccd936cf566d08e50d92caee136099d6df9cd603c3b7b18f33cb941b",
				Size:     226,
				Vsize:    226,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "22ba2a1c5c75040b84b6a2731e5d5c55dc3c57595caed950b218b3567ae8aa0d",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "30450220720793be4b7812a7b4b7ae6dd75cf37da09ddd4630debc09c34bf95782d16bdc022100f5f609629a559bcd96dcf82cb1657aa6374af78d3114c2f57d9f6fb204c97718[ALL] 03c584550e0b922323d39a357da5ee977e14a117fda1812ece91f33b06072cf49d", // nolint
							Hex: "4830450220720793be4b7812a7b4b7ae6dd75cf37da09ddd4630debc09c34bf95782d16bdc022100f5f609629a559bcd96dcf82cb1657aa6374af78d3114c2f57d9f6fb204c97718012103c584550e0b922323d39a357da5ee977e14a117fda1812ece91f33b06072cf49d", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 20.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 2deeef21adf2cf1d0fc8cf31efcf11af2321595d OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9142deeef21adf2cf1d0fc8cf31efcf11af2321595d88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D9KyBKj1shPhoXbkax3BVkttxfJh68A6W5",
							},
						},
					},
					{
						Value: 33.77810618,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 8a10bf7df695647d408a07d765699d41280663d8 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9148a10bf7df695647d408a07d765699d41280663d888ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DHj7rQMt57ZJTdpexHSkkK1ZMfLeS8RiXW",
							},
						},
					},
				},
			},
			{
				Hex:      "010000000830e60439eea77624e4b76424de72dc4345bad8e76d807459dce817155c39973c010000006b483045022021771a95a44f1acf88dcc0ac8b50f1756a96b648e0e8b130dfbc762c90d61806022100f40b7d1b643acf3b73ad32af99215fb8154cdd05b85d1d080fec154f7eefddb00121033779d94eb445e945e3c3ce7c815359e48a946dbf313b0375f979affd48d01fb4ffffffff5c506722bb545a465580685aeef32954da84574af3b805afbce40fac86029a5a010000006b483045022044b6c1506abf9823213034aea8c84f49435b1ac559bdc0bcd6d7d9f1d4b53207022100e22b75114b911a014cd42a24e243a42b22ee15422ef61cab61072ca547a9404e0121036774ad1017adb3fd13a756d7c76ad15c1dc49f37b71ea516c7cf93bbbd263c3cffffffff6afcf22f0f2db24119f4c3ddefee090be003430b72232170b8f3ab6c507074c6080000008c493046022100e36e614f3ec699fb3ecc9b7e496ddde2c91c583db9c6a64b94e46722d0e2dc22022100c55c3d38f915ecf046041b6e4594cddf86b27c5d83d7b291cb9815d58d33e7620141045d3e839d9eabede7593c0d89536f07904362dcbc2aa5cfeaa83aafa8ea9589a87bea418925a3db2f177b9c69d4fe48ed07b6960871e63828889e481a1e9c9a93ffffffffd0dd3a135633b1430afcd03e950dae8cacb097f0b8611b141124c057ddadca5a060000008c493046022100ec08de1e701770a39203c7063f8439e9f44fcb523b3c0aad33fa96d7a5e56988022100a36e1e438eb9b1856e9e03b7b783e1b1275045043f9aae10be42ab98e2428e6a0141045d3e839d9eabede7593c0d89536f07904362dcbc2aa5cfeaa83aafa8ea9589a87bea418925a3db2f177b9c69d4fe48ed07b6960871e63828889e481a1e9c9a93ffffffff9358deef5ff9ba99dc31189630f768cb4f879cb5764f587f71c1ccfd85ba51df150000008a47304402205536e7eee100668ec379b4fc39f1f33ada52b35514c0357d9df8a18750f675240220484100a77b44fd0c11ade5080e276558c3f9e23cc026539c86be1af30bd4e5720141045d3e839d9eabede7593c0d89536f07904362dcbc2aa5cfeaa83aafa8ea9589a87bea418925a3db2f177b9c69d4fe48ed07b6960871e63828889e481a1e9c9a93ffffffffd7f1a826d5cc65f19f91c0ec6fb84053e07683d6b115b4daf3b6a96548ddb0d2050000008b4830450220515a77a326e3e244b25cf9ba3baeb8636aad3bc2752a13c58d4792f65052c5b3022100b5edc14a0e2939dd3c0e352748513df88ca916aff9ea229cb0dd46f2d541ebb00141045d3e839d9eabede7593c0d89536f07904362dcbc2aa5cfeaa83aafa8ea9589a87bea418925a3db2f177b9c69d4fe48ed07b6960871e63828889e481a1e9c9a93ffffffff3cc94bda0348772bf4dd971f7f55a44f8131044fd92213d731dc2f5b4f2262b3000000006c4930460221008135f2ef5ffd16bb7d108cc71ebd23350933b479a73b9f48e2d682990fff5d29022100ac5a830a573e213dbbd4631ada472d523dcf64c072b58fbbef07613834430d93012102166b3359503760d953e25a996b1edac224b307d126392b2ac72a45105958891cffffffff5d6713f106f331dad761116e5550fa4ae73ef53a86e0542d9f98835cd9172d5e040000008a473044022006ebce6cb99658d714e226b2e941e490a32be41eb06fa8bc0ba824ea4cedf3c0022062bd83669800bd7ff0ac29b23fe5853cb61d633edb200af813a9d5cdea3e53da0141045d3e839d9eabede7593c0d89536f07904362dcbc2aa5cfeaa83aafa8ea9589a87bea418925a3db2f177b9c69d4fe48ed07b6960871e63828889e481a1e9c9a93ffffffff0200c96c4e170000001976a914e03b2b3d04e726d27a442d5421ec94decd16014888ac21371600000000001976a914803a518d40f3c7eea80fe8a1501a81f1d21074e488ac00000000", // nolint
				Hash:     "a63cc3784a264ac276d162ec2a3f3f609010e8852ceb3554a808a27fb6b228af",
				Size:     1423,
				Vsize:    1423,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "3c97395c1517e8dc5974806de7d8ba4543dc72de2464b7e42476a7ee3904e630",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "3045022021771a95a44f1acf88dcc0ac8b50f1756a96b648e0e8b130dfbc762c90d61806022100f40b7d1b643acf3b73ad32af99215fb8154cdd05b85d1d080fec154f7eefddb0[ALL] 033779d94eb445e945e3c3ce7c815359e48a946dbf313b0375f979affd48d01fb4", // nolint
							Hex: "483045022021771a95a44f1acf88dcc0ac8b50f1756a96b648e0e8b130dfbc762c90d61806022100f40b7d1b643acf3b73ad32af99215fb8154cdd05b85d1d080fec154f7eefddb00121033779d94eb445e945e3c3ce7c815359e48a946dbf313b0375f979affd48d01fb4", // nolint
						},
						Sequence: 4294967295,
					},

					{
						TxHash: "5a9a0286ac0fe4bcaf05b8f34a5784da5429f3ee5a688055465a54bb2267505c",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "3045022044b6c1506abf9823213034aea8c84f49435b1ac559bdc0bcd6d7d9f1d4b53207022100e22b75114b911a014cd42a24e243a42b22ee15422ef61cab61072ca547a9404e[ALL] 036774ad1017adb3fd13a756d7c76ad15c1dc49f37b71ea516c7cf93bbbd263c3c", // nolint
							Hex: "483045022044b6c1506abf9823213034aea8c84f49435b1ac559bdc0bcd6d7d9f1d4b53207022100e22b75114b911a014cd42a24e243a42b22ee15422ef61cab61072ca547a9404e0121036774ad1017adb3fd13a756d7c76ad15c1dc49f37b71ea516c7cf93bbbd263c3c", // nolint
						},
						Sequence: 4294967295,
					},

					{
						TxHash: "c67470506cabf3b8702123720b4303e00b09eeefddc3f41941b22d0f2ff2fc6a",
						Vout:   8,
						ScriptSig: &ScriptSig{
							ASM: "3046022100e36e614f3ec699fb3ecc9b7e496ddde2c91c583db9c6a64b94e46722d0e2dc22022100c55c3d38f915ecf046041b6e4594cddf86b27c5d83d7b291cb9815d58d33e762[ALL] 045d3e839d9eabede7593c0d89536f07904362dcbc2aa5cfeaa83aafa8ea9589a87bea418925a3db2f177b9c69d4fe48ed07b6960871e63828889e481a1e9c9a93", // nolint
							Hex: "493046022100e36e614f3ec699fb3ecc9b7e496ddde2c91c583db9c6a64b94e46722d0e2dc22022100c55c3d38f915ecf046041b6e4594cddf86b27c5d83d7b291cb9815d58d33e7620141045d3e839d9eabede7593c0d89536f07904362dcbc2aa5cfeaa83aafa8ea9589a87bea418925a3db2f177b9c69d4fe48ed07b6960871e63828889e481a1e9c9a93", // nolint
						},
						Sequence: 4294967295,
					},

					{
						TxHash: "5acaaddd57c02411141b61b8f097b0ac8cae0d953ed0fc0a43b13356133addd0",
						Vout:   6,
						ScriptSig: &ScriptSig{
							ASM: "3046022100ec08de1e701770a39203c7063f8439e9f44fcb523b3c0aad33fa96d7a5e56988022100a36e1e438eb9b1856e9e03b7b783e1b1275045043f9aae10be42ab98e2428e6a[ALL] 045d3e839d9eabede7593c0d89536f07904362dcbc2aa5cfeaa83aafa8ea9589a87bea418925a3db2f177b9c69d4fe48ed07b6960871e63828889e481a1e9c9a93", // nolint
							Hex: "493046022100ec08de1e701770a39203c7063f8439e9f44fcb523b3c0aad33fa96d7a5e56988022100a36e1e438eb9b1856e9e03b7b783e1b1275045043f9aae10be42ab98e2428e6a0141045d3e839d9eabede7593c0d89536f07904362dcbc2aa5cfeaa83aafa8ea9589a87bea418925a3db2f177b9c69d4fe48ed07b6960871e63828889e481a1e9c9a93", // nolint
						},
						Sequence: 4294967295,
					},

					{
						TxHash: "df51ba85fdccc1717f584f76b59c874fcb68f730961831dc99baf95fefde5893",
						Vout:   21,
						ScriptSig: &ScriptSig{
							ASM: "304402205536e7eee100668ec379b4fc39f1f33ada52b35514c0357d9df8a18750f675240220484100a77b44fd0c11ade5080e276558c3f9e23cc026539c86be1af30bd4e572[ALL] 045d3e839d9eabede7593c0d89536f07904362dcbc2aa5cfeaa83aafa8ea9589a87bea418925a3db2f177b9c69d4fe48ed07b6960871e63828889e481a1e9c9a93", // nolint
							Hex: "47304402205536e7eee100668ec379b4fc39f1f33ada52b35514c0357d9df8a18750f675240220484100a77b44fd0c11ade5080e276558c3f9e23cc026539c86be1af30bd4e5720141045d3e839d9eabede7593c0d89536f07904362dcbc2aa5cfeaa83aafa8ea9589a87bea418925a3db2f177b9c69d4fe48ed07b6960871e63828889e481a1e9c9a93", // nolint
						},
						Sequence: 4294967295,
					},

					{
						TxHash: "d2b0dd4865a9b6f3dab415b1d68376e05340b86fecc0919ff165ccd526a8f1d7",
						Vout:   5,
						ScriptSig: &ScriptSig{
							ASM: "30450220515a77a326e3e244b25cf9ba3baeb8636aad3bc2752a13c58d4792f65052c5b3022100b5edc14a0e2939dd3c0e352748513df88ca916aff9ea229cb0dd46f2d541ebb0[ALL] 045d3e839d9eabede7593c0d89536f07904362dcbc2aa5cfeaa83aafa8ea9589a87bea418925a3db2f177b9c69d4fe48ed07b6960871e63828889e481a1e9c9a93", // nolint
							Hex: "4830450220515a77a326e3e244b25cf9ba3baeb8636aad3bc2752a13c58d4792f65052c5b3022100b5edc14a0e2939dd3c0e352748513df88ca916aff9ea229cb0dd46f2d541ebb00141045d3e839d9eabede7593c0d89536f07904362dcbc2aa5cfeaa83aafa8ea9589a87bea418925a3db2f177b9c69d4fe48ed07b6960871e63828889e481a1e9c9a93", // nolint
						},
						Sequence: 4294967295,
					},

					{
						TxHash: "b362224f5b2fdc31d71322d94f0431814fa4557f1f97ddf42b774803da4bc93c",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "30460221008135f2ef5ffd16bb7d108cc71ebd23350933b479a73b9f48e2d682990fff5d29022100ac5a830a573e213dbbd4631ada472d523dcf64c072b58fbbef07613834430d93[ALL] 02166b3359503760d953e25a996b1edac224b307d126392b2ac72a45105958891c", // nolint
							Hex: "4930460221008135f2ef5ffd16bb7d108cc71ebd23350933b479a73b9f48e2d682990fff5d29022100ac5a830a573e213dbbd4631ada472d523dcf64c072b58fbbef07613834430d93012102166b3359503760d953e25a996b1edac224b307d126392b2ac72a45105958891c", // nolint
						},
						Sequence: 4294967295,
					},

					{
						TxHash: "5e2d17d95c83989f2d54e0863af53ee74afa50556e1161d7da31f306f113675d",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3044022006ebce6cb99658d714e226b2e941e490a32be41eb06fa8bc0ba824ea4cedf3c0022062bd83669800bd7ff0ac29b23fe5853cb61d633edb200af813a9d5cdea3e53da[ALL] 045d3e839d9eabede7593c0d89536f07904362dcbc2aa5cfeaa83aafa8ea9589a87bea418925a3db2f177b9c69d4fe48ed07b6960871e63828889e481a1e9c9a93", // nolint
							Hex: "473044022006ebce6cb99658d714e226b2e941e490a32be41eb06fa8bc0ba824ea4cedf3c0022062bd83669800bd7ff0ac29b23fe5853cb61d633edb200af813a9d5cdea3e53da0141045d3e839d9eabede7593c0d89536f07904362dcbc2aa5cfeaa83aafa8ea9589a87bea418925a3db2f177b9c69d4fe48ed07b6960871e63828889e481a1e9c9a93", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 1001.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 e03b2b3d04e726d27a442d5421ec94decd160148 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914e03b2b3d04e726d27a442d5421ec94decd16014888ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DRainManMRNWvKcuj5mQZKbsK4L2CTbAxt",
							},
						},
					},
					{
						Value: 0.01455905,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 803a518d40f3c7eea80fe8a1501a81f1d21074e4 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914803a518d40f3c7eea80fe8a1501a81f1d21074e488ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DGq6trgejnxZUDQDyPV52yM8CCXEbmHX1Y",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000024b88e1b58dfa49925b1ca091e7d1b8c43d546abff36c9f41e062aa2ec8b3b22d010000006c493046022100fca35a551d811f047d7bc8d2c96c48fca7425fdc56c23aa88f1a5d3e8dd590ec022100a7ff992c9eb3fbea30de1f98444b3855287c608fbe549e847d295c24c461235e0121030182db728d46cb7e09bac5ea81607f412b71d960e9337a905ff27ca46cfb2f6bffffffffeb2c6f0507d80ed63213c0dd699926f300c87c4e61aa7d075206b1c0b3dead8c000000006c493046022100b7a2a3e62f84521b76f0229009327b6d532805d865871a184955a3b608c8ef6a022100d64a49c81dd3063f86ea82190074bc33953ec8fb40b22267bca89e01578710cd012103bb3381cb1421d7c77193a25b80c978c8badc3a4bc4a9d6fa4450aae0ff77daa3ffffffff0100e87648170000001976a9140c5c994d75cc4c188c11dbf406a041f08344727688ac00000000", // nolint
				Hash:     "08d8a0f672d00310f079217852cefe84b25ff1a2a73a654b5df7066295751eba",
				Size:     342,
				Vsize:    342,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "2db2b3c82eaa62e0419f6cf3bf6a543dc4b8d1e791a01c5b9249fa8db5e1884b",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "3046022100fca35a551d811f047d7bc8d2c96c48fca7425fdc56c23aa88f1a5d3e8dd590ec022100a7ff992c9eb3fbea30de1f98444b3855287c608fbe549e847d295c24c461235e[ALL] 030182db728d46cb7e09bac5ea81607f412b71d960e9337a905ff27ca46cfb2f6b", // nolint
							Hex: "493046022100fca35a551d811f047d7bc8d2c96c48fca7425fdc56c23aa88f1a5d3e8dd590ec022100a7ff992c9eb3fbea30de1f98444b3855287c608fbe549e847d295c24c461235e0121030182db728d46cb7e09bac5ea81607f412b71d960e9337a905ff27ca46cfb2f6b", // nolint
						},
						Sequence: 4294967295,
					},
					{
						TxHash: "8caddeb3c0b10652077daa614e7cc800f3269969ddc01332d60ed807056f2ceb",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3046022100b7a2a3e62f84521b76f0229009327b6d532805d865871a184955a3b608c8ef6a022100d64a49c81dd3063f86ea82190074bc33953ec8fb40b22267bca89e01578710cd[ALL] 03bb3381cb1421d7c77193a25b80c978c8badc3a4bc4a9d6fa4450aae0ff77daa3", // nolint
							Hex: "493046022100b7a2a3e62f84521b76f0229009327b6d532805d865871a184955a3b608c8ef6a022100d64a49c81dd3063f86ea82190074bc33953ec8fb40b22267bca89e01578710cd012103bb3381cb1421d7c77193a25b80c978c8badc3a4bc4a9d6fa4450aae0ff77daa3", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 1000.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 0c5c994d75cc4c188c11dbf406a041f083447276 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9140c5c994d75cc4c188c11dbf406a041f08344727688ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D6GTaSW1p6afDyWGUfeCpooPtDF6kessdh",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000002edef21f2e6dbd1da2ee5371bdfbe9290359002cbe6816486a520873be909e6a6000000006b483045022100b0c043dc01eeb9dcf3851e8b41a101a01626b0c0d9ed051acc1822fb44b7818c022042752baafbfea2d8311577aa1c96582280904130c9df864c3428551ef6e4454e012102184363f5cfdd6138a5e574dba3fba87583d0073252e23542e572f284dc29d25fffffffff8e869d76156f8b545ccb31ff5ec09b4746bf6412a0a1e1c64d0f956a1aa851f7010000006a47304402204d1520941c1eb5cbd1f21a7d3a583c6ee240600faeda6ef5f881ba7a80c519cc02204e327651ef12d337176a324d52aa8a83b9511fa0dbab436c9995a5fec2e7778a012102184363f5cfdd6138a5e574dba3fba87583d0073252e23542e572f284dc29d25fffffffff0200205fa0120000001976a914aa72ef734a5414fc4249fcc4aee8a3b430b2892e88ac00eadd8d040000001976a91472285613b3fd6a077445991cad0046969ba5a62788ac00000000", // nolint
				Hash:     "000fa0f6c12c4770cc8ea70912dea0519e79333fe042d9c3978450ada52e06a4",
				Size:     373,
				Vsize:    373,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "a6e609e93b8720a5866481e6cb0290359092bedf1b37e52edad1dbe6f221efed",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3045022100b0c043dc01eeb9dcf3851e8b41a101a01626b0c0d9ed051acc1822fb44b7818c022042752baafbfea2d8311577aa1c96582280904130c9df864c3428551ef6e4454e[ALL] 02184363f5cfdd6138a5e574dba3fba87583d0073252e23542e572f284dc29d25f", // nolint
							Hex: "483045022100b0c043dc01eeb9dcf3851e8b41a101a01626b0c0d9ed051acc1822fb44b7818c022042752baafbfea2d8311577aa1c96582280904130c9df864c3428551ef6e4454e012102184363f5cfdd6138a5e574dba3fba87583d0073252e23542e572f284dc29d25f", // nolint
						},
						Sequence: 4294967295,
					},
					{
						TxHash: "f751a81a6a950f4dc6e1a1a01264bf46479bc05eff31cb5c548b6f15769d868e",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "304402204d1520941c1eb5cbd1f21a7d3a583c6ee240600faeda6ef5f881ba7a80c519cc02204e327651ef12d337176a324d52aa8a83b9511fa0dbab436c9995a5fec2e7778a[ALL] 02184363f5cfdd6138a5e574dba3fba87583d0073252e23542e572f284dc29d25f", // nolint
							Hex: "47304402204d1520941c1eb5cbd1f21a7d3a583c6ee240600faeda6ef5f881ba7a80c519cc02204e327651ef12d337176a324d52aa8a83b9511fa0dbab436c9995a5fec2e7778a012102184363f5cfdd6138a5e574dba3fba87583d0073252e23542e572f284dc29d25f", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 800.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 aa72ef734a5414fc4249fcc4aee8a3b430b2892e OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914aa72ef734a5414fc4249fcc4aee8a3b430b2892e88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DLgM6ibPHqhFjjT7niZ6mWHxFfGCqNoqVS",
							},
						},
					},
					{
						Value: 195.60000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 72285613b3fd6a077445991cad0046969ba5a627 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a91472285613b3fd6a077445991cad0046969ba5a62788ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DFYhtjm7L9aLGPiCmv8UJVxdU7dUd9iZa4",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000001ab6d5cb4ff11e8cdd710c403a8d33e5c3eafd87d7185879818b063f7efb8548e010000006a4730440220554a6519fe2755db4c1a3310574cafc0601cde7a424c898d86b43d1732a6b566022038ec16a38b8f7c69effe12c2c1eff4a29895ed446d173b32c65759b8161f09e0012103b8c8e44691a1811f355dc610fd322028e1f6937a2a5d32097e891dae605b810cffffffff0200e40b54020000001976a9143824aa2f78626ad32d3d6fc23ec0bf1380df4bfe88ac008dbce3180000001976a914c24e84b87e6ff62318799e2d86969b1cb3a5a68f88ac00000000", // nolint
				Hash:     "d8eb74127bf13a43161fabe39d989e48abeaec956244645eb0da2d2fcb0b408b",
				Size:     225,
				Vsize:    225,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "8e54b8eff763b018988785717dd8af3e5c3ed3a803c410d7cde811ffb45c6dab",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "30440220554a6519fe2755db4c1a3310574cafc0601cde7a424c898d86b43d1732a6b566022038ec16a38b8f7c69effe12c2c1eff4a29895ed446d173b32c65759b8161f09e0[ALL] 03b8c8e44691a1811f355dc610fd322028e1f6937a2a5d32097e891dae605b810c", // nolint
							Hex: "4730440220554a6519fe2755db4c1a3310574cafc0601cde7a424c898d86b43d1732a6b566022038ec16a38b8f7c69effe12c2c1eff4a29895ed446d173b32c65759b8161f09e0012103b8c8e44691a1811f355dc610fd322028e1f6937a2a5d32097e891dae605b810c", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 100.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 3824aa2f78626ad32d3d6fc23ec0bf1380df4bfe OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9143824aa2f78626ad32d3d6fc23ec0bf1380df4bfe88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DAFxJWph4KVuTjm1czTpq1kT27eZbug7t2",
							},
						},
					},
					{
						Value: 1069.00000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 c24e84b87e6ff62318799e2d86969b1cb3a5a68f OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914c24e84b87e6ff62318799e2d86969b1cb3a5a68f88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DNrVh2iY47eWdswyKTXgeXV2cEPUSjPTmZ",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000001cdbee0abc61024d2d9d40e7119d73d5770203d7415b7297811a7e9e102ac5520000000006a4730440220648eb9b8160073b298fe4f5e3ce044c6b78665bd739fe44defc618ccfc839ec602201fa37293e3a466e746b690c9688bf97a6cc74572dcd4ab632750705574939025012102a8a7f85c202ea64ac08fd248388e60f62688d0505ea7210d7a36cd6d00ac3648ffffffff0236d96447810f00001976a914bbacd0ef6d81ba9910e2c8428041a65442a370ea88acc5973993240000001976a914ac0b37041f4e2b6337aee944364a614c942ac4fc88ac00000000", // nolint
				Hash:     "841934645bc2cfa78eee0f9a037836de7f1eb8f608fdedac41e8b4835ef26bb3",
				Size:     225,
				Vsize:    225,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "2055ac02e1e9a7117829b715743d2070573dd719710ed4d9d22410c6abe0becd",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "30440220648eb9b8160073b298fe4f5e3ce044c6b78665bd739fe44defc618ccfc839ec602201fa37293e3a466e746b690c9688bf97a6cc74572dcd4ab632750705574939025[ALL] 02a8a7f85c202ea64ac08fd248388e60f62688d0505ea7210d7a36cd6d00ac3648", // nolint
							Hex: "4730440220648eb9b8160073b298fe4f5e3ce044c6b78665bd739fe44defc618ccfc839ec602201fa37293e3a466e746b690c9688bf97a6cc74572dcd4ab632750705574939025012102a8a7f85c202ea64ac08fd248388e60f62688d0505ea7210d7a36cd6d00ac3648", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 170479.22989366,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 bbacd0ef6d81ba9910e2c8428041a65442a370ea OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914bbacd0ef6d81ba9910e2c8428041a65442a370ea88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DNFRvcC7pn4sjMsFxCBxRSDRzRKkEJmXoh",
							},
						},
					},
					{
						Value: 1570.88847813,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 ac0b37041f4e2b6337aee944364a614c942ac4fc OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914ac0b37041f4e2b6337aee944364a614c942ac4fc88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DLpnCUiRLSCzTH1bJvERhnPKJz24esh12K",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000010591e2ab8b8a503aa69bfd5688862cee828d43c14c09bcd78e3920ba2fa37f18010000006a473044022079b1d0f7f649b8e5e8634f7deeb6482a57b0af090770e061b2c8d1cfba71cc3002203ce0560db5aca86716971b11fb9b03723d90272a49376d78ef0ecee7dff810e40121038fa4eabe1d7463d13d6c86a92fcc3e8000473adc41cb07c372ba9e06a3f43194ffffffff0280c52d8b010000001976a9143c4d481aa7776cc867de67143b9f9df65a445ff488ac8021f416030000001976a914706631056aa360576922550a808c268cfa578cd988ac00000000", // nolint
				Hash:     "edeac6b327f5baa9013f149135f0a55b07c98ceffcfba8cd07be30bf0eceab75",
				Size:     225,
				Vsize:    225,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "187fa32fba20398ed7bc094cc1438d82ee2c868856fd9ba63a508a8babe29105",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3044022079b1d0f7f649b8e5e8634f7deeb6482a57b0af090770e061b2c8d1cfba71cc3002203ce0560db5aca86716971b11fb9b03723d90272a49376d78ef0ecee7dff810e4[ALL] 038fa4eabe1d7463d13d6c86a92fcc3e8000473adc41cb07c372ba9e06a3f43194", // nolint
							Hex: "473044022079b1d0f7f649b8e5e8634f7deeb6482a57b0af090770e061b2c8d1cfba71cc3002203ce0560db5aca86716971b11fb9b03723d90272a49376d78ef0ecee7dff810e40121038fa4eabe1d7463d13d6c86a92fcc3e8000473adc41cb07c372ba9e06a3f43194", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 66.30000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 3c4d481aa7776cc867de67143b9f9df65a445ff4 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9143c4d481aa7776cc867de67143b9f9df65a445ff488ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DAdwfPgdkUxnFtbHKL3dga1p5nM7uFBuso",
							},
						},
					},
					{
						Value: 132.70000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 706631056aa360576922550a808c268cfa578cd9 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914706631056aa360576922550a808c268cfa578cd988ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DFPQe8eEKsReCgFhXKSFYeZnFpwVqJ8H76",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000013ac16d8d789c616a9908d1b2d8405c49594dd76f85b0102fda9cdaa2600f69e9010000006b483045022100e86e4db2ecd8e13333c6c6ea53576f9f66c19c7ed9bc3e4e62ad37a3fc89caa4022057d05333637ea5b9d76df538ce4302f0798f6e5d1e48d49eaba5bc629c1d41f2012103d7c35f0c7c24c510ba334ae5f3f7d86e5c0af94b8d2ccd7adc052b2f7e1d6b4dffffffff0200b1bcdd040000001976a914c3fa20522ee60c30e8e90f0d508336d49655823088ac00943577000000001976a914ae3d035ba09bce98918bb13021e650f1dfb3b08588ac00000000", // nolint
				Hash:     "ec54cae5e048cce6310dc5260ad2beb37475887f1552da559ecfe19aa1927f1f",
				Size:     226,
				Vsize:    226,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "0b50ae7e1ff3f3c63c44e4f5580d24fda1716cfceca6945f4301c652b662b316",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "3045022100e86e4db2ecd8e13333c6c6ea53576f9f66c19c7ed9bc3e4e62ad37a3fc89caa4022057d05333637ea5b9d76df538ce4302f0798f6e5d1e48d49eaba5bc629c1d41f2[ALL] 03d7c35f0c7c24c510ba334ae5f3f7d86e5c0af94b8d2ccd7adc052b2f7e1d6b4d", // nolint
							Hex: "483045022100e86e4db2ecd8e13333c6c6ea53576f9f66c19c7ed9bc3e4e62ad37a3fc89caa4022057d05333637ea5b9d76df538ce4302f0798f6e5d1e48d49eaba5bc629c1d41f2012103d7c35f0c7c24c510ba334ae5f3f7d86e5c0af94b8d2ccd7adc052b2f7e1d6b4d", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 209.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 c3fa20522ee60c30e8e90f0d508336d496558230 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914c3fa20522ee60c30e8e90f0d508336d49655823088ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DP1Kwjzw78c44NVBEi46uHkU2cVsHpLiyJ",
							},
						},
					},
					{
						Value: 20.00000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 ae3d035ba09bce98918bb13021e650f1dfb3b085 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914ae3d035ba09bce98918bb13021e650f1dfb3b08588ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DM2PCuJe7k9Qw3N4e5aHpBPChCiprK5ZNt",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000001688cd393eaf5aea8516e11dc5f395a449bab1ca6f5c6fc679cac6bd781d5d85d4c0000006b48304502204ce49efb622982dba187c0767c163560fd272f032470af9b1f7c0de04120b783022100d108d2f2878c78d5a4afcb4a1259310752f648e5383e10bec24cae296a1ab4ad012103979d6d85a49d18deb678c3dbfd9ec73a82ee61971c08e63d2e4db37361094212ffffffff01bdbe798c320000001976a9147936b878100ccc012d4eecc7124accc6774078de88ac00000000", // nolint
				Hash:     "e5d1ea180aee5f2aad4d110cc3c166b4a1ae9f8af12905cbca4b6dc93bc8d1a7",
				Size:     192,
				Vsize:    192,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "5dd8d581d76bac9c67fcc6f5a61cab9b445a395fdc116e51a8aef5ea93d38c68",
						Vout:   76,
						ScriptSig: &ScriptSig{
							ASM: "304502204ce49efb622982dba187c0767c163560fd272f032470af9b1f7c0de04120b783022100d108d2f2878c78d5a4afcb4a1259310752f648e5383e10bec24cae296a1ab4ad[ALL] 03979d6d85a49d18deb678c3dbfd9ec73a82ee61971c08e63d2e4db37361094212", // nolint
							Hex: "48304502204ce49efb622982dba187c0767c163560fd272f032470af9b1f7c0de04120b783022100d108d2f2878c78d5a4afcb4a1259310752f648e5383e10bec24cae296a1ab4ad012103979d6d85a49d18deb678c3dbfd9ec73a82ee61971c08e63d2e4db37361094212", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 2171.05153725,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 7936b878100ccc012d4eecc7124accc6774078de OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9147936b878100ccc012d4eecc7124accc6774078de88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DGC1rXGKhU7GMwYZF4TSFbEMyCoSwZSeAr",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000019cb1fb0be32828d9a5fb4d8d2cad34e1d351a1a5990692c09fc9522f3c788d52010000006a47304402206a3fbe36f6e9c09ca3e6920cc5ec214950bec7db3144c407c8d58f0c7b1f05830220585881d965f59e53e9e5b83c3617c10595e7d362ebf7e00c0ba8bb3c2c102e100121023134cc9962bd55c25349192e30686f7e06b1e71effb83ea98bc142ab71d4976affffffff0100e721a2040000001976a914d68cb846739bd68d4bfe6346ca0e7de32268af8e88ac00000000", // nolint
				Hash:     "208120e984209bfff10ca2cb1ef646df252c61eb72c7acae77da980a2001d240",
				Size:     191,
				Vsize:    191,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "528d783c2f52c99fc0920699a5a151d3e134ad2c8d4dfba5d92828e30bfbb19c",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "304402206a3fbe36f6e9c09ca3e6920cc5ec214950bec7db3144c407c8d58f0c7b1f05830220585881d965f59e53e9e5b83c3617c10595e7d362ebf7e00c0ba8bb3c2c102e10[ALL] 023134cc9962bd55c25349192e30686f7e06b1e71effb83ea98bc142ab71d4976a", // nolint
							Hex: "47304402206a3fbe36f6e9c09ca3e6920cc5ec214950bec7db3144c407c8d58f0c7b1f05830220585881d965f59e53e9e5b83c3617c10595e7d362ebf7e00c0ba8bb3c2c102e100121023134cc9962bd55c25349192e30686f7e06b1e71effb83ea98bc142ab71d4976a", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 199.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 d68cb846739bd68d4bfe6346ca0e7de32268af8e OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914d68cb846739bd68d4bfe6346ca0e7de32268af8e88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DQhXihGETDz5Antxbx8eymhnJFcUBhiTg2",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000001a42ca3dce24910f351f385912d8a37ea87cee8917ad4a90a95ccccd193405e13000000006b483045022100e1636331779607619dab8254b599f7191d5c9a1139529f38a8a7aefc406a152402206ae7739b822c205ae044507da67b6e356688841d093ec87d76cc60f7d223ff5001210318c5bd2030e074abc39f5b34f3642500dd2f33eb686f78983fb639417ad0645bffffffff0220b5a664310200001976a91453645d0ee901a62b5b91590c6cda95e42409cbe888ac80f0fa02000000001976a91488a65118c4d13c1a526f4de4fff41b9b35d816cc88ac00000000", // nolint
				Hash:     "8941e4ad43afc0093761b2b23ba712a44cd2b530db9a8c251761e2d681724281",
				Size:     226,
				Vsize:    226,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "135e4093d1cccc950aa9d47a91e8ce87ea378a2d9185f351f31049e2dca32ca4",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3045022100e1636331779607619dab8254b599f7191d5c9a1139529f38a8a7aefc406a152402206ae7739b822c205ae044507da67b6e356688841d093ec87d76cc60f7d223ff50[ALL] 0318c5bd2030e074abc39f5b34f3642500dd2f33eb686f78983fb639417ad0645b", // nolint
							Hex: "483045022100e1636331779607619dab8254b599f7191d5c9a1139529f38a8a7aefc406a152402206ae7739b822c205ae044507da67b6e356688841d093ec87d76cc60f7d223ff5001210318c5bd2030e074abc39f5b34f3642500dd2f33eb686f78983fb639417ad0645b", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 24111.65300000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 53645d0ee901a62b5b91590c6cda95e42409cbe8 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a91453645d0ee901a62b5b91590c6cda95e42409cbe888ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DCk2rZPRBZkkKjrCjrqRdLiujEyaYTLCvK",
							},
						},
					},
					{
						Value: 0.50000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 88a65118c4d13c1a526f4de4fff41b9b35d816cc OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a91488a65118c4d13c1a526f4de4fff41b9b35d816cc88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DHbdgEt42PiKmAzbJpHfypEeygaGZmcay3",
							},
						},
					},
				},
			},
			{
				Hex:      "010000000184b219412d00a2e9431f60de222b1a96f4edf2f4c886bf7d5c2e226f4cf2f367000000006a47304402200aea45983a1b2fd6d3c521b9dc8b9b7a90f11f68c4beec7eacb755125c1074c5022054a938bfc624d5b7c0f74bf4eff2620babf41739675786f8f95f9ea2f89ca15901210210afe85ead8d32a4f338cc9b5e3ebe7e68fccfa1b05ceb846d51174bb6e5c8b8ffffffff0200cce06ccc0000001976a91485648c7c54342b76befc69a7f6a02d0b7e3c047988ac00e1f505000000001976a914289888fd9f4d032d43aad4485d5dba6b931834e988ac00000000", // nolint
				Hash:     "cd12c012b713348073c65955047eeffc83aa27953e90c017ecb64b38345b483d",
				Size:     225,
				Vsize:    225,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "67f3f24c6f222e5c7dbf86c8f4f2edf4961a2b22de601f43e9a2002d4119b284",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "304402200aea45983a1b2fd6d3c521b9dc8b9b7a90f11f68c4beec7eacb755125c1074c5022054a938bfc624d5b7c0f74bf4eff2620babf41739675786f8f95f9ea2f89ca159[ALL] 0210afe85ead8d32a4f338cc9b5e3ebe7e68fccfa1b05ceb846d51174bb6e5c8b8", // nolint
							Hex: "47304402200aea45983a1b2fd6d3c521b9dc8b9b7a90f11f68c4beec7eacb755125c1074c5022054a938bfc624d5b7c0f74bf4eff2620babf41739675786f8f95f9ea2f89ca15901210210afe85ead8d32a4f338cc9b5e3ebe7e68fccfa1b05ceb846d51174bb6e5c8b8", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 8780.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 85648c7c54342b76befc69a7f6a02d0b7e3c0479 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a91485648c7c54342b76befc69a7f6a02d0b7e3c047988ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DHJQs2HniggCDFPJjqc5SS13bCy17Eznwa",
							},
						},
					},
					{
						Value: 1.00000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 289888fd9f4d032d43aad4485d5dba6b931834e9 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914289888fd9f4d032d43aad4485d5dba6b931834e988ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D8qkJFto7FVUFTVhCATJLQrHikHYmuxM6K",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000001de633d4e614214b2c4f630f82f78cbf01cdbdb31265f99480204149a0d736dbf000000006a473044022068c5edbae3e6626bc7da3029d0df0cd59114ad5d2e8533f4f26439251361d88302204a940ea94bebb8e9d1d93822fa23e11c23037c7971f36b3b3e839aa243d4e2c3012103a93d999bb383a6bd9cf802101b66245e2c6cbbce302ac964893a8e7f187e9461ffffffff0200a3e111000000001976a9148849b263f7d808f2fc97241944b94f2849b2601488ac00b5cab32b0000001976a914a6f534602558d53a1aa3ed98bb8816ba3ceac02088ac00000000", // nolint
				Hash:     "8b65f57e9e275c4ae531ddd5f4c4e93c295e5ec877ae847c4df55cd43ba155ab",
				Size:     225,
				Vsize:    225,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "bf6d730d9a14040248995f2631dbdb1cf0cb782ff830f6c4b21442614e3d63de",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3044022068c5edbae3e6626bc7da3029d0df0cd59114ad5d2e8533f4f26439251361d88302204a940ea94bebb8e9d1d93822fa23e11c23037c7971f36b3b3e839aa243d4e2c3[ALL] 03a93d999bb383a6bd9cf802101b66245e2c6cbbce302ac964893a8e7f187e9461", // nolint
							Hex: "473044022068c5edbae3e6626bc7da3029d0df0cd59114ad5d2e8533f4f26439251361d88302204a940ea94bebb8e9d1d93822fa23e11c23037c7971f36b3b3e839aa243d4e2c3012103a93d999bb383a6bd9cf802101b66245e2c6cbbce302ac964893a8e7f187e9461", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 3.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 8849b263f7d808f2fc97241944b94f2849b26014 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9148849b263f7d808f2fc97241944b94f2849b2601488ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DHZiitmo22bYapKrh3H3jMPps4kyEXMK8p",
							},
						},
					},
					{
						Value: 1877.00000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 a6f534602558d53a1aa3ed98bb8816ba3ceac020 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914a6f534602558d53a1aa3ed98bb8816ba3ceac02088ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DLMtTDUy5gSeSvLyu6hEzJ4giwV3xQ5Bdq",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000001d9b3c3b27f52d72e93c88e2bdf54ddfdb27027119401eb0c33a4b17b0640e2c0000000006b48304502203747cb4c101b4b2b8d3815cccd9baeefb36cea2f595e6306db56b242a08f74c2022100a8f83b573154c23f82f8d1ca002fb10fe8c55033ae86625d6a1644a9b00f27fd01210316d1f584e2a315261fefb1b1c77cdc22eecc4e34244ff46742100083a4a62498ffffffff0282fabeab080000001976a91412ece6e9180a15251177f66e39f0dca7fe1947c688acf3ba0401000000001976a9140be8cc678017ea9c9eafeaac63ddfa31fc2136e088ac00000000", // nolint
				Hash:     "3e9aa63e33df82b71856d9182fe4119e8ff7e203f183588f99ca0ac1c7c3051d",
				Size:     226,
				Vsize:    226,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "c0e240067bb1a4330ceb0194112770b2fddd54df2b8ec8932ed7527fb2c3b3d9",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "304502203747cb4c101b4b2b8d3815cccd9baeefb36cea2f595e6306db56b242a08f74c2022100a8f83b573154c23f82f8d1ca002fb10fe8c55033ae86625d6a1644a9b00f27fd[ALL] 0316d1f584e2a315261fefb1b1c77cdc22eecc4e34244ff46742100083a4a62498", // nolint
							Hex: "48304502203747cb4c101b4b2b8d3815cccd9baeefb36cea2f595e6306db56b242a08f74c2022100a8f83b573154c23f82f8d1ca002fb10fe8c55033ae86625d6a1644a9b00f27fd01210316d1f584e2a315261fefb1b1c77cdc22eecc4e34244ff46742100083a4a62498", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 372.41158274,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 12ece6e9180a15251177f66e39f0dca7fe1947c6 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a91412ece6e9180a15251177f66e39f0dca7fe1947c688ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D6sAVx5nFxo9SUkkFDMK7hmLuw1tnGJfo5",
							},
						},
					},
					{
						Value: 0.17087219,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 0be8cc678017ea9c9eafeaac63ddfa31fc2136e0 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9140be8cc678017ea9c9eafeaac63ddfa31fc2136e088ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D6E4rVJQ5XG5Bzh7sXenSTdRfD8KuyDp8t",
							},
						},
					},
				},
			},
			{
				Hex:      "010000000168423f1633c4f503cc0971638c89e00f453dff809b11aefe65e48db83869f348000000006b48304502206dd4ecf7e1c79cb139b8bf4507f48f1bd5a3ce207ec10807496ca7a7dabf4526022100ba0a3b99eba4e0cb86d175fe82873b165b2f06f65663f43929424e6632d561290121021bfcc1615687c72e64cc9d06f57501602a09a3e8e612cdc9a66bee2668f3a3a1ffffffff0261c7bcae070000001976a914a24be5ff2adcb5fc0c5e360b553ceb22a425a64b88acbdbd174c130000001976a914d690e204b63c98c0470b2b7ea3367ba109719aa088ac00000000", // nolint
				Hash:     "3f7508e5ff526b142031c26d5bae369621bd065051271366cf7f9bdb90dd94e1",
				Size:     226,
				Vsize:    226,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "48f36938b88de465feae119b80ff3d450fe0898c637109cc03f5c433163f4268",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "304502206dd4ecf7e1c79cb139b8bf4507f48f1bd5a3ce207ec10807496ca7a7dabf4526022100ba0a3b99eba4e0cb86d175fe82873b165b2f06f65663f43929424e6632d56129[ALL] 021bfcc1615687c72e64cc9d06f57501602a09a3e8e612cdc9a66bee2668f3a3a1", // nolint
							Hex: "48304502206dd4ecf7e1c79cb139b8bf4507f48f1bd5a3ce207ec10807496ca7a7dabf4526022100ba0a3b99eba4e0cb86d175fe82873b165b2f06f65663f43929424e6632d561290121021bfcc1615687c72e64cc9d06f57501602a09a3e8e612cdc9a66bee2668f3a3a1", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 329.96378465,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 a24be5ff2adcb5fc0c5e360b553ceb22a425a64b OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914a24be5ff2adcb5fc0c5e360b553ceb22a425a64b88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DKwEvrFZdZNZRWH74MDhQjn6NeM262bPpj",
							},
						},
					},
					{
						Value: 828.81002941,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 d690e204b63c98c0470b2b7ea3367ba109719aa0 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914d690e204b63c98c0470b2b7ea3367ba109719aa088ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DQhchx2AFu5dhRacN5QY5gL4EK2z1TAwHs",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000001f12e57c8b3456b9dc8a7403938a2822a57843fdb537543a73332d7e5e78f85f7000000006b48304502206978be3154ae94d58f6abd4cd916bce8f587d53f1646868463edbaaca82ab6c7022100cdd208ea5bfcaabfe37801a5fd29cd83db2e2ae28cf4a229fa1fa58beadf90860121036a59c0bad2fcda23a69c78adea325b58466ca59d5126bcc69faea043beaa1e25ffffffff014a2a95c60e0000001976a914e2c4c95cdb866f579f93977d668de345b8388e6a88ac00000000", // nolint
				Hash:     "b454b575dd861b20f73747ac3e34b9b333551e262c815923ad460390d911857a",
				Size:     192,
				Vsize:    192,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "f7858fe7e5d73233a7437553db3f84572a82a2383940a7c89d6b45b3c8572ef1",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "304502206978be3154ae94d58f6abd4cd916bce8f587d53f1646868463edbaaca82ab6c7022100cdd208ea5bfcaabfe37801a5fd29cd83db2e2ae28cf4a229fa1fa58beadf9086[ALL] 036a59c0bad2fcda23a69c78adea325b58466ca59d5126bcc69faea043beaa1e25", // nolint
							Hex: "48304502206978be3154ae94d58f6abd4cd916bce8f587d53f1646868463edbaaca82ab6c7022100cdd208ea5bfcaabfe37801a5fd29cd83db2e2ae28cf4a229fa1fa58beadf90860121036a59c0bad2fcda23a69c78adea325b58466ca59d5126bcc69faea043beaa1e25", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 634.61206602,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 e2c4c95cdb866f579f93977d668de345b8388e6a OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914e2c4c95cdb866f579f93977d668de345b8388e6a88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DRp8zamCGLSYRv2W8dEnFBFbYfmEi6VpHk",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000002acd75dffa180f9fa929653841e9c61c8553d067d00ea9a7141dad6880984a63f010000006a47304402206513136d876906f9917998e344cb080473e441c0be2b59900dae0221b04d077a02203409c91aae44643cdd98da8597372d53edd019af6936974eb04b292d53848b27012102a9bc6c4bc078f805cd1c7aca4404069adc7717879ebdd679ee1a246c7e6f4bf3ffffffffd19f37347a4549891e3dfd0e7eab911d9e67196ddcbf41502ffef1773bd7a98b000000006b48304502203d3f097f0eeea4f998347d842f564dc9e3f5d0612b4dbff3bfb84298814484e302210086e60299707e98f069aa99f9f05e1494b3ccae4e1cdd1e62c36dbf6a6049d27f012102e8d9f3ea4da6cb2dd65250fe694c94bf27788a98d5ef6b02e62da5f0bb6ab20bffffffff021d29473a000000001976a914af0f58fc983699ce76cb6ef18fc271dcdfb79ff588ac00904e5a030000001976a914fb1a23206fad89723d855e5a80ccb9daa149244088ac00000000", // nolint
				Hash:     "b67637b5f675c49ad7731665bb82ab974b845a9e6d3eeb8ea1f6e9578f013b83",
				Size:     373,
				Vsize:    373,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "3fa6840988d6da41719aea007d063d55c8619c1e84539692faf980a1ff5dd7ac",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "304402206513136d876906f9917998e344cb080473e441c0be2b59900dae0221b04d077a02203409c91aae44643cdd98da8597372d53edd019af6936974eb04b292d53848b27[ALL] 02a9bc6c4bc078f805cd1c7aca4404069adc7717879ebdd679ee1a246c7e6f4bf3", // nolint
							Hex: "47304402206513136d876906f9917998e344cb080473e441c0be2b59900dae0221b04d077a02203409c91aae44643cdd98da8597372d53edd019af6936974eb04b292d53848b27012102a9bc6c4bc078f805cd1c7aca4404069adc7717879ebdd679ee1a246c7e6f4bf3", // nolint
						},
						Sequence: 4294967295,
					},
					{
						TxHash: "8ba9d73b77f1fe2f5041bfdc6d19679e1d91ab7e0efd3d1e8949457a34379fd1",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "304502203d3f097f0eeea4f998347d842f564dc9e3f5d0612b4dbff3bfb84298814484e302210086e60299707e98f069aa99f9f05e1494b3ccae4e1cdd1e62c36dbf6a6049d27f[ALL] 02e8d9f3ea4da6cb2dd65250fe694c94bf27788a98d5ef6b02e62da5f0bb6ab20b", // nolint
							Hex: "48304502203d3f097f0eeea4f998347d842f564dc9e3f5d0612b4dbff3bfb84298814484e302210086e60299707e98f069aa99f9f05e1494b3ccae4e1cdd1e62c36dbf6a6049d27f012102e8d9f3ea4da6cb2dd65250fe694c94bf27788a98d5ef6b02e62da5f0bb6ab20b", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 9.77742109,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 af0f58fc983699ce76cb6ef18fc271dcdfb79ff5 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914af0f58fc983699ce76cb6ef18fc271dcdfb79ff588ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DM6jBCVhdnJvAaeBzfsdixjp1uypGvPQnd",
							},
						},
					},
					{
						Value: 144.00000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 fb1a23206fad89723d855e5a80ccb9daa1492440 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914fb1a23206fad89723d855e5a80ccb9daa149244088ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DU2oTTFDgYyJyLSjhMU6991HRZbJb3VDMU",
							},
						},
					},
				},
			},
			{
				Hex:      "010000000133a2e5e307a65f0345ea33901d307f63000139b3132a288d9a66016a8bb0abb8010000006a47304402200c66f408e869b1506f117c993d009df3b075a68b9ce4e62f1c92f783dcd860d702204b6aefac1539e90dd35b9c3fb0c8256a2f3c11c77c98c65121bfb80deba1f96b012102f7740849c2c6c4be311a1fb669e15b8af78bd201d15a4433d82d981724e3dfe8ffffffff0200eec1ee150000001976a914515ff6aef100b80aa0bc0f825e4197d398b4885688ac006d7c4d000000001976a914016b5977c4a98b00ce4bdd731d42896d3947149b88ac00000000", // nolint
				Hash:     "48fc3402cbe9e71e0e8a6aeaece9bd252312639e95db6cd4d555324ac491e003",
				Size:     225,
				Vsize:    225,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "b8abb08b6a01669a8d282a13b3390100637f301d9033ea45035fa607e3e5a233",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "304402200c66f408e869b1506f117c993d009df3b075a68b9ce4e62f1c92f783dcd860d702204b6aefac1539e90dd35b9c3fb0c8256a2f3c11c77c98c65121bfb80deba1f96b[ALL] 02f7740849c2c6c4be311a1fb669e15b8af78bd201d15a4433d82d981724e3dfe8", // nolint
							Hex: "47304402200c66f408e869b1506f117c993d009df3b075a68b9ce4e62f1c92f783dcd860d702204b6aefac1539e90dd35b9c3fb0c8256a2f3c11c77c98c65121bfb80deba1f96b012102f7740849c2c6c4be311a1fb669e15b8af78bd201d15a4433d82d981724e3dfe8", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 942.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 515ff6aef100b80aa0bc0f825e4197d398b48856 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914515ff6aef100b80aa0bc0f825e4197d398b4885688ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DCZNETvXEdgq7mntj56FqyZn1uSgqi641X",
							},
						},
					},
					{
						Value: 13.00000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 016b5977c4a98b00ce4bdd731d42896d3947149b OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914016b5977c4a98b00ce4bdd731d42896d3947149b88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D5GbpDtyU8dGvJeLf1Hh7FiUFRsZk3Ttuo",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000012371105379e58d8fa4d403ccd3d0bbd2f360e49e49d6fd81d2b94ca9aeba2290000000006c4930460221008210281a12b168f02c3f45e79996ff06f3aa7cfa25cf956d54c9bbc59403a7c30221008ac2008fe26f403d43e254c9ff4ef9534229bb7137aa57184e7e415354fa030f01210362027b5a4959973f0c005653e7bfbcdb557099ee42e9053aabf6138594fb2ebeffffffff0200ff0f270b0000001976a9148e6da23aea89b44790320e2821d77735b86d711288ac006d7c4d000000001976a91428dae9efa934162c0c4e17ea75687731eff374f188ac00000000", // nolint
				Hash:     "65e32f36d5e2d42dfe54219e610882eee13b9102ba345ea271f5bf4872ee2f72",
				Size:     227,
				Vsize:    227,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "9022baaea94cb9d281fdd6499ee460f3d2bbd0d3cc03d4a48f8de57953107123",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "30460221008210281a12b168f02c3f45e79996ff06f3aa7cfa25cf956d54c9bbc59403a7c30221008ac2008fe26f403d43e254c9ff4ef9534229bb7137aa57184e7e415354fa030f[ALL] 0362027b5a4959973f0c005653e7bfbcdb557099ee42e9053aabf6138594fb2ebe", // nolint
							Hex: "4930460221008210281a12b168f02c3f45e79996ff06f3aa7cfa25cf956d54c9bbc59403a7c30221008ac2008fe26f403d43e254c9ff4ef9534229bb7137aa57184e7e415354fa030f01210362027b5a4959973f0c005653e7bfbcdb557099ee42e9053aabf6138594fb2ebe", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 479.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 8e6da23aea89b44790320e2821d77735b86d7112 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9148e6da23aea89b44790320e2821d77735b86d711288ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DJ8BpyvJ8zFAcYd2NdN1pPMVwBBW8QK8LC",
							},
						},
					},
					{
						Value: 13.00000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 28dae9efa934162c0c4e17ea75687731eff374f1 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a91428dae9efa934162c0c4e17ea75687731eff374f188ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D8s7pL6YtCcDrXbcjgnp6c7N94LW31MQHn",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000001d0b5959c1487e2d4f3de202b52ce9286cba3060e3a8ebce06c735472d67bc813000000006c4930460221009fcf3285a07360f3fabd1306f66b1ac965e9d178ca9ceaf7487a292d96cc8792022100d74380ddb3da48f986ffb8b063124a97a79a1500f5bfd5181d97fb0cb026c7a3012102487c3f3cc9a5f443107c36fa06844b64dd8b23896d9491f6a6d698bdb511c8e9ffffffff02f00d7bcb010000001976a914318eced404c11e3c5fde840d5aaeeb9efc00c6b688ac005ed0b2000000001976a91462e9340893d0cf62a9e14eeed711e5764386f05e88ac00000000", // nolint
				Hash:     "85ea9de43dfd49716b7908f227c72dc99c9c82bc54da8ff1e3cc488d4a4ac80c",
				Size:     227,
				Vsize:    227,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "13c87bd67254736ce0bc8e3a0e06a3cb8692ce522b20def3d4e287149c95b5d0",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "30460221009fcf3285a07360f3fabd1306f66b1ac965e9d178ca9ceaf7487a292d96cc8792022100d74380ddb3da48f986ffb8b063124a97a79a1500f5bfd5181d97fb0cb026c7a3[ALL] 02487c3f3cc9a5f443107c36fa06844b64dd8b23896d9491f6a6d698bdb511c8e9", // nolint
							Hex: "4930460221009fcf3285a07360f3fabd1306f66b1ac965e9d178ca9ceaf7487a292d96cc8792022100d74380ddb3da48f986ffb8b063124a97a79a1500f5bfd5181d97fb0cb026c7a3012102487c3f3cc9a5f443107c36fa06844b64dd8b23896d9491f6a6d698bdb511c8e9", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 77.08806640,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 318eced404c11e3c5fde840d5aaeeb9efc00c6b6 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914318eced404c11e3c5fde840d5aaeeb9efc00c6b688ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D9f8j8Hxi1dNrhKtNZtoY8uDc2awvqZ1mQ",
							},
						},
					},
					{
						Value: 30.00000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 62e9340893d0cf62a9e14eeed711e5764386f05e OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a91462e9340893d0cf62a9e14eeed711e5764386f05e88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DEA68JbFSzcSeGn5bGe7Xf6w5XbBdbuHkg",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000001c79ae575877f2670a88a53b8939c10ab6413c1ec6847f2381b9dd7d44f69bdde000000006c493046022100ed5787095d60824035f10043d2c218cffd7022f149bc6d0ad6aad7cc8a1ff7ba0221009a0fd8a6b9b74a2fba9b16da60061d5c44fd58b7e279ddc0bc5b604ada381d8901210369f29d6f5caac452d06eccb2899fdcd03ce385a050ab32c868ff29f7397eff68ffffffff0298bad24f040000001976a914171ff4c844a9eef6bb79cf374e684a70158e7d5388aca93bb11f050000001976a914a7124f6adeccce0e17e94f9a68161c52a32eafa288ac00000000", // nolint
				Hash:     "f774f5b7d2b95fee611314bd859005cc3c73fb5f47d5e60f84dd9f6b2a972fe9",
				Size:     227,
				Vsize:    227,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "debd694fd4d79d1b38f24768ecc11364ab109c93b8538aa870267f8775e59ac7",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3046022100ed5787095d60824035f10043d2c218cffd7022f149bc6d0ad6aad7cc8a1ff7ba0221009a0fd8a6b9b74a2fba9b16da60061d5c44fd58b7e279ddc0bc5b604ada381d89[ALL] 0369f29d6f5caac452d06eccb2899fdcd03ce385a050ab32c868ff29f7397eff68", // nolint
							Hex: "493046022100ed5787095d60824035f10043d2c218cffd7022f149bc6d0ad6aad7cc8a1ff7ba0221009a0fd8a6b9b74a2fba9b16da60061d5c44fd58b7e279ddc0bc5b604ada381d8901210369f29d6f5caac452d06eccb2899fdcd03ce385a050ab32c868ff29f7397eff68", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 185.19079576,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 171ff4c844a9eef6bb79cf374e684a70158e7d53 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914171ff4c844a9eef6bb79cf374e684a70158e7d5388ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D7FNN2Q3uZYgVQcbtAN5PYoQUvD7Z4gHAW",
							},
						},
					},
					{
						Value: 220.06545321,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 a7124f6adeccce0e17e94f9a68161c52a32eafa2 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914a7124f6adeccce0e17e94f9a68161c52a32eafa288ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DLNVKWaeVmjpWQrPgcvYY4YaiuQD6SKVAQ",
							},
						},
					},
				},
			},
			{
				Hex:      "010000000135a2a0611e0b326f05e48124e7a7592c0ca52a7f9a3efcadf03d74027a244018000000006c493046022100e296ff2f5a1b996c922c327c371b4e21f2fb40bdbba2f9b5befb5c74b09b9f8d02210091f07e275b035a784e4ed0a5c13c7836a90b52b739f2b3efcf47fab0dd6de59d0121038569c188cbc829f851ccda57a3235617c1909e1de2419c9c1169529b14fe0c00ffffffff0200d2496b000000001976a914b80be6e1a7bdd4b984bd5680db5179f7b645031488ac00410b38080000001976a914fc32581d33bd5eacda6158ea10cf4fd9c635d84d88ac00000000", // nolint
				Hash:     "02b1d923d46274cba9f4937b75eee3348896f409d11edcff3613f0f55a38a2a1",
				Size:     227,
				Vsize:    227,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "1840247a02743df0adfc3e9a7f2aa50c2c59a7e72481e4056f320b1e61a0a235",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3046022100e296ff2f5a1b996c922c327c371b4e21f2fb40bdbba2f9b5befb5c74b09b9f8d02210091f07e275b035a784e4ed0a5c13c7836a90b52b739f2b3efcf47fab0dd6de59d[ALL] 038569c188cbc829f851ccda57a3235617c1909e1de2419c9c1169529b14fe0c00", // nolint
							Hex: "493046022100e296ff2f5a1b996c922c327c371b4e21f2fb40bdbba2f9b5befb5c74b09b9f8d02210091f07e275b035a784e4ed0a5c13c7836a90b52b739f2b3efcf47fab0dd6de59d0121038569c188cbc829f851ccda57a3235617c1909e1de2419c9c1169529b14fe0c00", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 18.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 b80be6e1a7bdd4b984bd5680db5179f7b6450314 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914b80be6e1a7bdd4b984bd5680db5179f7b645031488ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DMvF8WfzRPPbt1DhhyQh5YYyckzWfyNTLQ",
							},
						},
					},
					{
						Value: 353.00000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 fc32581d33bd5eacda6158ea10cf4fd9c635d84d OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914fc32581d33bd5eacda6158ea10cf4fd9c635d84d88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DU8b8aM38855uAvTzJeAdjMb5b236uqDUi",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000038055caf093c8a044fa5429adfff3484e117fe4fe76cb27174d2b9d792b4fbb62010000006b48304502205ed4304d094741126acec996aff3fbca1c3ceac1d931886c5b9eb2575361df04022100a563593827d8f1e4d10f9cf325f51b901fc450adf0b2ff584ca4965ed1dbbb2d0121025916328cd6c38bf5f98693323158a4325fa8b333a215c2ae0b1ea4a92f673a38ffffffff8a1a976bcd51d0cf22b021b9803d7cbeb77a33e1a2dfc6ca4773320dcce7c8f7000000006b483045022100d446953828ddb670baa707ccb1ef9937d3e48660a1c90891f40ed983f71eaa0d02200413f37a60fd0cff0a99cd1e0a14d37c839e460c56771753d6f8fc99a6939927012102d7ddada145d1d5c081427f3913ce9f77e849425a2f2c4a4b66737b22a1c7488cffffffff7607d8257d5b492b8988e58186c00b6d214a401ae14105b7d9e89983fd33f584000000006c493046022100cf6bc9a14f631a69895662fad0e809e7afe4d1c11ed2c1348c40bcbc0aa6fe95022100c8fac9fbee9d8d6f637fcdcfb27c3aaa962bd80e9819bc5d64fceb078842bcd7012102d7ddada145d1d5c081427f3913ce9f77e849425a2f2c4a4b66737b22a1c7488cffffffff0280e06200000000001976a914841377b2b993ed53a7333b09856f1e02bc791ebd88ac00a3e111000000001976a914ab1b5cae43c1b925e0b89c2aa78494f1bc6f924e88ac00000000", // nolint
				Hash:     "80d242bec4ccd84ac083d912230eb439abda5f8c53e15b6504836530506f372d",
				Size:     523,
				Vsize:    523,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "62bb4f2b799d2b4d1727cb76fee47f114e48f3ffad2954fa44a0c893f0ca5580",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "304502205ed4304d094741126acec996aff3fbca1c3ceac1d931886c5b9eb2575361df04022100a563593827d8f1e4d10f9cf325f51b901fc450adf0b2ff584ca4965ed1dbbb2d[ALL] 025916328cd6c38bf5f98693323158a4325fa8b333a215c2ae0b1ea4a92f673a38", // nolint
							Hex: "48304502205ed4304d094741126acec996aff3fbca1c3ceac1d931886c5b9eb2575361df04022100a563593827d8f1e4d10f9cf325f51b901fc450adf0b2ff584ca4965ed1dbbb2d0121025916328cd6c38bf5f98693323158a4325fa8b333a215c2ae0b1ea4a92f673a38", // nolint
						},
						Sequence: 4294967295,
					},
					{
						TxHash: "f7c8e7cc0d327347cac6dfa2e1337ab7be7c3d80b921b022cfd051cd6b971a8a",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3045022100d446953828ddb670baa707ccb1ef9937d3e48660a1c90891f40ed983f71eaa0d02200413f37a60fd0cff0a99cd1e0a14d37c839e460c56771753d6f8fc99a6939927[ALL] 02d7ddada145d1d5c081427f3913ce9f77e849425a2f2c4a4b66737b22a1c7488c", // nolint
							Hex: "483045022100d446953828ddb670baa707ccb1ef9937d3e48660a1c90891f40ed983f71eaa0d02200413f37a60fd0cff0a99cd1e0a14d37c839e460c56771753d6f8fc99a6939927012102d7ddada145d1d5c081427f3913ce9f77e849425a2f2c4a4b66737b22a1c7488c", // nolint
						},
						Sequence: 4294967295,
					},
					{
						TxHash: "84f533fd8399e8d9b70541e11a404a216d0bc08681e588892b495b7d25d80776",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3046022100cf6bc9a14f631a69895662fad0e809e7afe4d1c11ed2c1348c40bcbc0aa6fe95022100c8fac9fbee9d8d6f637fcdcfb27c3aaa962bd80e9819bc5d64fceb078842bcd7[ALL] 02d7ddada145d1d5c081427f3913ce9f77e849425a2f2c4a4b66737b22a1c7488c", // nolint
							Hex: "493046022100cf6bc9a14f631a69895662fad0e809e7afe4d1c11ed2c1348c40bcbc0aa6fe95022100c8fac9fbee9d8d6f637fcdcfb27c3aaa962bd80e9819bc5d64fceb078842bcd7012102d7ddada145d1d5c081427f3913ce9f77e849425a2f2c4a4b66737b22a1c7488c", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 100000.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 101f0445d2cee10c1f820dac0fdab961c35c5940 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914101f0445d2cee10c1f820dac0fdab961c35c594088ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D6cLWQTCnvCxUZS5G99dsDmKc3vGuo26w6",
							},
						},
					},
					{
						Value: 0.06480000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 841377b2b993ed53a7333b09856f1e02bc791ebd OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914841377b2b993ed53a7333b09856f1e02bc791ebd88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DHBT4DAgnYXZEouEqT6dwXJxe1KpuLvGhw",
							},
						},
					},
					{
						Value: 3.00000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 ab1b5cae43c1b925e0b89c2aa78494f1bc6f924e OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914ab1b5cae43c1b925e0b89c2aa78494f1bc6f924e88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DLjpsD4z6GPqogtoEhTLgwkm9zM3xXQNu2",
							},
						},
					},
				},
			},
			{
				Hex:      "010000000181427281d6e26117258c9adb30b5d24ca412a73bb2b2613709c0af43ade44189000000006b483045022100d6bef6d56c92a90075e50ab20a0260e7e5fd5a4cec6d4943ef245bcb35df749d022013ad806ba85663476b7e0ee89ed0fa05ae2846932deeb620ff43e0d983f9a7ca012102c35c73eecb3f7498836706fab553a8a9881b5ab42ad1b726cac11a5e66f29e4bffffffff02a002c055310200001976a914605ab623319733a2ff3c7f749214cde3ab4cb05a88ac80f0fa02000000001976a914ed79e21df195f9701e91eebc2b4e69257058101b88ac00000000", // nolint
				Hash:     "6f7d46c9731bfa51c5bcd76ac6b40a1a8ef85a192973f02f66a82d8f261d500d",
				Size:     226,
				Vsize:    226,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "8941e4ad43afc0093761b2b23ba712a44cd2b530db9a8c251761e2d681724281",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3045022100d6bef6d56c92a90075e50ab20a0260e7e5fd5a4cec6d4943ef245bcb35df749d022013ad806ba85663476b7e0ee89ed0fa05ae2846932deeb620ff43e0d983f9a7ca[ALL] 02c35c73eecb3f7498836706fab553a8a9881b5ab42ad1b726cac11a5e66f29e4b", // nolint
							Hex: "483045022100d6bef6d56c92a90075e50ab20a0260e7e5fd5a4cec6d4943ef245bcb35df749d022013ad806ba85663476b7e0ee89ed0fa05ae2846932deeb620ff43e0d983f9a7ca012102c35c73eecb3f7498836706fab553a8a9881b5ab42ad1b726cac11a5e66f29e4b", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 24109.15300000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 605ab623319733a2ff3c7f749214cde3ab4cb05a OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914605ab623319733a2ff3c7f749214cde3ab4cb05a88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DDva5RC5e8Kg85dw83HCW9a7y9p2HuWxLX",
							},
						},
					},
					{
						Value: 0.50000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 ed79e21df195f9701e91eebc2b4e69257058101b OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914ed79e21df195f9701e91eebc2b4e69257058101b88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DSnkhKXeVz8Np58tVzMWrcGkUh4twmQKMf",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000016c873218e5dc062965e5a5179057cd27bf44c99abec60288c645f1c1c8ba9bfb000000006a47304402201412ad826f09485dcd32c11c521c80269fa077d1e19eefa828d40135e795b2bb02207dd7b8b4c065092e3f16038dec8676910aaa362d753416a3ef1600e24acca3b6012102d43a6e77148fa7584674817fccd77e564832bc623158b598ebde6812ef22b590ffffffff02a029c1e1010000001976a91445a8fb7e94982c104f1121d549c5f8d28df3ce8988ac603fa71d000000001976a9147de21cf2eb1b16f6b50dda20b9045d2ef729993c88ac00000000", // nolint
				Hash:     "0e4a7e7c42d7a991b5c491a9c5500ebf83fd1b266448a6dc986bea3a7e14c179",
				Size:     225,
				Vsize:    225,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "fb9bbac8c1f145c68802c6be9ac944bf27cd579017a5e5652906dce51832876c",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "304402201412ad826f09485dcd32c11c521c80269fa077d1e19eefa828d40135e795b2bb02207dd7b8b4c065092e3f16038dec8676910aaa362d753416a3ef1600e24acca3b6[ALL] 02d43a6e77148fa7584674817fccd77e564832bc623158b598ebde6812ef22b590", // nolint
							Hex: "47304402201412ad826f09485dcd32c11c521c80269fa077d1e19eefa828d40135e795b2bb02207dd7b8b4c065092e3f16038dec8676910aaa362d753416a3ef1600e24acca3b6012102d43a6e77148fa7584674817fccd77e564832bc623158b598ebde6812ef22b590", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 80.82500000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 45a8fb7e94982c104f1121d549c5f8d28df3ce89 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a91445a8fb7e94982c104f1121d549c5f8d28df3ce8988ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DBVRbbd43CC2A9KyWknvcU2xaKaCXwaFuH",
							},
						},
					},
					{
						Value: 4.97500000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 7de21cf2eb1b16f6b50dda20b9045d2ef729993c OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9147de21cf2eb1b16f6b50dda20b9045d2ef729993c88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DGchsrEEsxCp9yz2QTyJ2QyhypGzR3JGMR",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000001246964d523055dcbd0e1844643ee90b05adadd1805dc5a8138e1d572109f1eb6000000006a47304402200f4874872c3f96659b59c63e219c7827e131242538ee58910ef7d37c7a440bda02200509d82aa2738a383f1f76f495ce37a20cc56530cd1139aaa009c4e4700d5f720121035f33520e83f96e1dfa8136e987471c906568959fa40ba8c0ca73a576bbe03b8dffffffff024054adf4000000001976a9146e782f433c0ed10dec2156fcf750da0f60f84f6988ac603fa71d000000001976a9148238d5ea6fdcec5b76b8792205be25b559386ac888ac00000000", // nolint
				Hash:     "1791cf26e46774074e750ebf5248a75d570446e6777de4a1ae61c64e8f4f3638",
				Size:     225,
				Vsize:    225,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "b61e9f1072d5e138815adc0518ddda5ab090ee434684e1d0cb5d0523d5646924",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "304402200f4874872c3f96659b59c63e219c7827e131242538ee58910ef7d37c7a440bda02200509d82aa2738a383f1f76f495ce37a20cc56530cd1139aaa009c4e4700d5f72[ALL] 035f33520e83f96e1dfa8136e987471c906568959fa40ba8c0ca73a576bbe03b8d", // nolint
							Hex: "47304402200f4874872c3f96659b59c63e219c7827e131242538ee58910ef7d37c7a440bda02200509d82aa2738a383f1f76f495ce37a20cc56530cd1139aaa009c4e4700d5f720121035f33520e83f96e1dfa8136e987471c906568959fa40ba8c0ca73a576bbe03b8d", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 41.05000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 6e782f433c0ed10dec2156fcf750da0f60f84f69 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9146e782f433c0ed10dec2156fcf750da0f60f84f6988ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DFDCqw6jZMSdxdSBxNS7xd65Bz28KbQijQ",
							},
						},
					},
					{
						Value: 4.97500000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 8238d5ea6fdcec5b76b8792205be25b559386ac8 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9148238d5ea6fdcec5b76b8792205be25b559386ac888ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DH1eUCmWH43iSyyFAyrhYGQCbJ6eKbT1yf",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000001c067ae3dbd1d8112ac2eee4e0a9455a035ff8079a57bcf035359f01896509042000000006b48304502200764a38588e468ee488372e49b111c2c060ef1dbb37bc5165216c2de7a799fe7022100bf5111c8c8f6098440513399abc45305404ab70b572feaec84379a193251c984012102ab2e8ebc9a14e3a2fea3a5dd1949a9c075664c2c57ab1b360a3ced60b5010e40ffffffff026cfb8e15000000001976a9149a41d5bc2ab936b2f6de2483cb0bef50169b6d9988ac00ca9a3b000000001976a91455c71f61392272e15347c027382677cfe67f5ad788ac00000000", // nolint
				Hash:     "3b232ef20f087ac6e6045a799327a801f05034e29801c4b96d95b5767904f38c",
				Size:     226,
				Vsize:    226,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "4290509618f0595303cf7ba57980ff35a055940a4eee2eac12811dbd3dae67c0",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "304502200764a38588e468ee488372e49b111c2c060ef1dbb37bc5165216c2de7a799fe7022100bf5111c8c8f6098440513399abc45305404ab70b572feaec84379a193251c984[ALL] 02ab2e8ebc9a14e3a2fea3a5dd1949a9c075664c2c57ab1b360a3ced60b5010e40", // nolint
							Hex: "48304502200764a38588e468ee488372e49b111c2c060ef1dbb37bc5165216c2de7a799fe7022100bf5111c8c8f6098440513399abc45305404ab70b572feaec84379a193251c984012102ab2e8ebc9a14e3a2fea3a5dd1949a9c075664c2c57ab1b360a3ced60b5010e40", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 3.61692012,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 9a41d5bc2ab936b2f6de2483cb0bef50169b6d99 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9149a41d5bc2ab936b2f6de2483cb0bef50169b6d9988ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DKCjU69TmzEqdqUgDUQb1SYhe4tPKwz128",
							},
						},
					},
					{
						Value: 10.00000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 55c71f61392272e15347c027382677cfe67f5ad7 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a91455c71f61392272e15347c027382677cfe67f5ad788ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DCxeWqrvFEdTKUvuSVaZh5FpA1TGE9R3fF",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000001262ed6364d4b45d407650a5374a37699b6c5ae01a70cfe034abe33c981f24b93010000006b483045022064f3cb97cf6af3a9fbb4f528b772c4486ce1050e560bc459e9d8756cbf2749de022100bfd821d47570fcbd01209f69c6b75a75c754f23a15ae9f365af59e2b88fc4f45012102b4ede20a853e1391b4a29f72fa1a76fc519ecd3b561d11ea8c0abfc638ab93b3ffffffff02a6e28d24010000001976a9149a383ebe681d88e0314d1dfb4f596dd17aa422bc88ac00ca9a3b000000001976a9147cb9e6d120be2bd152b748f8f590c926b47bf9c788ac00000000", // nolint
				Hash:     "aae3a0dda5e481b11d8bd5c2c2107c54453dea6eaf084fefe8d5ff8ecec3ae48",
				Size:     226,
				Vsize:    226,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "934bf281c933be4a03fe0ca701aec5b69976a374530a6507d4454b4d36d62e26",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "3045022064f3cb97cf6af3a9fbb4f528b772c4486ce1050e560bc459e9d8756cbf2749de022100bfd821d47570fcbd01209f69c6b75a75c754f23a15ae9f365af59e2b88fc4f45[ALL] 02b4ede20a853e1391b4a29f72fa1a76fc519ecd3b561d11ea8c0abfc638ab93b3", // nolint
							Hex: "483045022064f3cb97cf6af3a9fbb4f528b772c4486ce1050e560bc459e9d8756cbf2749de022100bfd821d47570fcbd01209f69c6b75a75c754f23a15ae9f365af59e2b88fc4f45012102b4ede20a853e1391b4a29f72fa1a76fc519ecd3b561d11ea8c0abfc638ab93b3", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 49.08245670,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 9a383ebe681d88e0314d1dfb4f596dd17aa422bc OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9149a383ebe681d88e0314d1dfb4f596dd17aa422bc88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DKCXyn553nzqsB2Q9CcxGA1rcdb8aKcZ64",
							},
						},
					},
					{
						Value: 10.00000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 7cb9e6d120be2bd152b748f8f590c926b47bf9c7 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9147cb9e6d120be2bd152b748f8f590c926b47bf9c788ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DGWb2iYGJ58zm7ZYmjNkbpuJmJRVPUCQw8",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000010ee2653a9516e70405dccde6c2155514970fb6aa0d72b4145c844fbd929127ea000000006b4830450221009240281b701818e78f67bff775204071a9ea3215512acd8e7fff8886938cf257022043bb000937fd5c0cd99fbba19936731ec0b9a0b31c92742f62eaba2bd44625d9012102a296f5b4fd4e83cd78b5dc71e328b6ad97205f77137e26fa9c918f50ca29e71effffffff020064488a000000001976a91471f4bb8d47f4a0c0d0d87d091abb7405861e892488ac603fa71d000000001976a9141cd3b22d4ebb6f7d94187e4f7fedcea24cee11c688ac00000000", // nolint
				Hash:     "62486d4e4f6ddd5cf24e90562781e374c4b656a60cc5152242027a13283db7ae",
				Size:     226,
				Vsize:    226,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "ea279192bd4f845c14b4720daab60f97145515c2e6cddc0504e716953a65e20e",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "30450221009240281b701818e78f67bff775204071a9ea3215512acd8e7fff8886938cf257022043bb000937fd5c0cd99fbba19936731ec0b9a0b31c92742f62eaba2bd44625d9[ALL] 02a296f5b4fd4e83cd78b5dc71e328b6ad97205f77137e26fa9c918f50ca29e71e", // nolint
							Hex: "4830450221009240281b701818e78f67bff775204071a9ea3215512acd8e7fff8886938cf257022043bb000937fd5c0cd99fbba19936731ec0b9a0b31c92742f62eaba2bd44625d9012102a296f5b4fd4e83cd78b5dc71e328b6ad97205f77137e26fa9c918f50ca29e71e", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 23.20000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 71f4bb8d47f4a0c0d0d87d091abb7405861e8924 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a91471f4bb8d47f4a0c0d0d87d091abb7405861e892488ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DFXe5FoteTdcZYkxA7gWtXHzwQTXDM9f33",
							},
						},
					},
					{
						Value: 4.97500000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 1cd3b22d4ebb6f7d94187e4f7fedcea24cee11c6 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9141cd3b22d4ebb6f7d94187e4f7fedcea24cee11c688ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D7mX4XTPV5TDYqvqvaaXf9QCAdxyaHf45C",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000016d0c87bfbb6b049498ab2ab01f1bf492cf56daca99da882415c418d6700e1dd2010000006b48304502210096eb511d7c04a47fb8e2165db508253f7af87de4378a6830eb622a3527684df102203d9d823dea7ff0af1526dbb031f0dccbf520f6b294dd0c3993e05427f1c2cc2b01210214c0c53161cd86b81a01c8b12321a5f1fe8ab93b7b9ea0186fc7d606bbd03870ffffffff0200b33f71000000001976a914434608760e385a4b8d54178d41631b201c804e5b88ac008c8647000000001976a91445716aa96fbab4d30fe0355254a0954368d0987e88ac00000000", // nolint
				Hash:     "0634e7ba2bfb7ccdf1ce40958488b401ca5aa2127029c0e4a790703cbe2bd771",
				Size:     226,
				Vsize:    226,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "d21d0e70d618c4152488da99cada56cf92f41b1fb02aab9894046bbbbf870c6d",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "304502210096eb511d7c04a47fb8e2165db508253f7af87de4378a6830eb622a3527684df102203d9d823dea7ff0af1526dbb031f0dccbf520f6b294dd0c3993e05427f1c2cc2b[ALL] 0214c0c53161cd86b81a01c8b12321a5f1fe8ab93b7b9ea0186fc7d606bbd03870", // nolint
							Hex: "48304502210096eb511d7c04a47fb8e2165db508253f7af87de4378a6830eb622a3527684df102203d9d823dea7ff0af1526dbb031f0dccbf520f6b294dd0c3993e05427f1c2cc2b01210214c0c53161cd86b81a01c8b12321a5f1fe8ab93b7b9ea0186fc7d606bbd03870", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 19.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 434608760e385a4b8d54178d41631b201c804e5b OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914434608760e385a4b8d54178d41631b201c804e5b88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DBGoi6LQvCPHBypL2odYiqQphiFtNknr8f",
							},
						},
					},
					{
						Value: 12.00000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 45716aa96fbab4d30fe0355254a0954368d0987e OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a91445716aa96fbab4d30fe0355254a0954368d0987e88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DBUH2pc4s4GBPAmWDPQ7phrGDBnFrAAZ2K",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000011d05c3c7c10aca998f5883f103e2f78f9e11e42f18d95618b782df333ea69a3e000000006b4830450221009b51194a564c121966444c476cf92b64c4667a152cb0c06875c01608924bfee702205b90d6837dd8814fcc2928d9d472521db25d5feaee69ee03e40852b75a6d57cf0121023732f8acf3eef3fea2226f2de4802b8d789d2b30909d4e9009c9d5138591605cffffffff027d9fd81c060000001976a914a5c1bad6b2e517d6df2ec67bae14e9ba77b5244f88ac057af088020000001976a9144bb19b49d47e77c216554eb718012ee1a8d58b7688ac00000000", // nolint
				Hash:     "140e621560e9af41de49460ae139d6dc0599c1d0e83b55fa9abe49415f683d1a",
				Size:     226,
				Vsize:    226,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "3e9aa63e33df82b71856d9182fe4119e8ff7e203f183588f99ca0ac1c7c3051d",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "30450221009b51194a564c121966444c476cf92b64c4667a152cb0c06875c01608924bfee702205b90d6837dd8814fcc2928d9d472521db25d5feaee69ee03e40852b75a6d57cf[ALL] 023732f8acf3eef3fea2226f2de4802b8d789d2b30909d4e9009c9d5138591605c", // nolint
							Hex: "4830450221009b51194a564c121966444c476cf92b64c4667a152cb0c06875c01608924bfee702205b90d6837dd8814fcc2928d9d472521db25d5feaee69ee03e40852b75a6d57cf0121023732f8acf3eef3fea2226f2de4802b8d789d2b30909d4e9009c9d5138591605c", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 262.53762429,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 a5c1bad6b2e517d6df2ec67bae14e9ba77b5244f OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914a5c1bad6b2e517d6df2ec67bae14e9ba77b5244f88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DLFY7VhXCE2Kw7YdZEySH7Jamxp7kiGY9q",
							},
						},
					},
					{
						Value: 108.87395845,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 4bb19b49d47e77c216554eb718012ee1a8d58b76 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9144bb19b49d47e77c216554eb718012ee1a8d58b7688ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DC3Kyy7dFLcRSVEpNQEQdfL5XjLnKFqJGm",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000001e3789e9ae0d8271cde731c58f6ef02ccace4df32efa1d7aa9f719b07751a97e6010000006c493046022100c82f89daba81ffbfd05c5cb85a8238c7950e0ff6ef79994249552910eb77fee70221009f0e2c840e809d7414262c6dcb651cdaf030ff12c47b122b4409e82ef13823b401210389a12c3c8443c73ccbed2e2d529bea8edaf7a4efb8fa6c5b40d3ef68591f4ab5ffffffff0200fab459010000001976a9147ee2e69aa7226562024951cb09b15ae0520a06f888ac00943577000000001976a9144db9be540350646dc34fa78def8e04fed49d81ce88ac00000000", // nolint
				Hash:     "682205324a6040743acc93bc79023ed3f2d5e13ac671ed60f3b20ed4d6e6ec6f",
				Size:     227,
				Vsize:    227,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "e6971a75079b719faad7a1ef32dfe4accc02eff6581c73de1c27d8e09a9e78e3",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "3046022100c82f89daba81ffbfd05c5cb85a8238c7950e0ff6ef79994249552910eb77fee70221009f0e2c840e809d7414262c6dcb651cdaf030ff12c47b122b4409e82ef13823b4[ALL] 0389a12c3c8443c73ccbed2e2d529bea8edaf7a4efb8fa6c5b40d3ef68591f4ab5", // nolint
							Hex: "493046022100c82f89daba81ffbfd05c5cb85a8238c7950e0ff6ef79994249552910eb77fee70221009f0e2c840e809d7414262c6dcb651cdaf030ff12c47b122b4409e82ef13823b401210389a12c3c8443c73ccbed2e2d529bea8edaf7a4efb8fa6c5b40d3ef68591f4ab5", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 58.00000000,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 7ee2e69aa7226562024951cb09b15ae0520a06f8 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9147ee2e69aa7226562024951cb09b15ae0520a06f888ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DGi1Vmn8vCXwYzr2QXcPmLEfcL6Mu2Xt2h",
							},
						},
					},
					{
						Value: 20.00000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 4db9be540350646dc34fa78def8e04fed49d81ce OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9144db9be540350646dc34fa78def8e04fed49d81ce88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DCE55iF3wTpAZjqdtddSvaaA2PgixJjSUG",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000001b36bf25e83b4e841acedfd08f6b81e7fde3678039a0fee8ea7cfc25b64341984000000006c493046022100a5ba6a649289c4bc5f21bc0bc0bba4b30ccee57197aead9eef44244c180bd6be022100cbb2a297103c85c434288bdd3196233231a763750e347f2bac212a10343e057e012102f4bb5a222354a9649740a8b8b300f3e1a0c7ade18a409bae7c2ea9f0ecac0575ffffffff027e04ed7ca70e00001976a914c660254bafe140f4a6ec8c248b61a5b7a427ba6788acb8f381c4d90000001976a914e0a032f63f59e7b7b9a79b9dcd7bb8121f672f3688ac00000000", // nolint
				Hash:     "5385158b8721c90f47d1dd1e05aec4076f1b614c5b5b64950b8c556c8e4b6213",
				Size:     227,
				Vsize:    227,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "841934645bc2cfa78eee0f9a037836de7f1eb8f608fdedac41e8b4835ef26bb3",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3046022100a5ba6a649289c4bc5f21bc0bc0bba4b30ccee57197aead9eef44244c180bd6be022100cbb2a297103c85c434288bdd3196233231a763750e347f2bac212a10343e057e[ALL] 02f4bb5a222354a9649740a8b8b300f3e1a0c7ade18a409bae7c2ea9f0ecac0575", // nolint
							Hex: "493046022100a5ba6a649289c4bc5f21bc0bc0bba4b30ccee57197aead9eef44244c180bd6be022100cbb2a297103c85c434288bdd3196233231a763750e347f2bac212a10343e057e012102f4bb5a222354a9649740a8b8b300f3e1a0c7ade18a409bae7c2ea9f0ecac0575", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 161125.18235262,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 c660254bafe140f4a6ec8c248b61a5b7a427ba67 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914c660254bafe140f4a6ec8c248b61a5b7a427ba6788ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DPE1WZSw5YbyK2VS3QYjpauyu7xTdMztjC",
							},
						},
					},
					{
						Value: 9353.04754104,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 e0a032f63f59e7b7b9a79b9dcd7bb8121f672f36 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914e0a032f63f59e7b7b9a79b9dcd7bb8121f672f3688ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DRcop49Apg1Bt9zZrzFg3c3i1jKFWvBeE7",
							},
						},
					},
				},
			},
			{
				Hex:      "010000000113624b8e6c558c0b95645b5b4c611b6f07c4ae051eddd1470fc921878b158553000000006b483045022100ccdfdb7f9b45e5bf1d36b4963fa05714b0cb963e3f93d370f806166203904cd9022031d5d036bd92716e742756511a4be0ca85674f1d65cf8285eccb20bec9c09873012103d3e092b233a629befb6a2249024faebfddf944140f5403a7eedb3f011c47271dffffffff020c987287ec0d00001976a914c5207f84f009632513b304af1e69757f79d1df4588ac728b84efba0000001976a9143d9e306694eb9a55e00113a6f9b838392d15c71a88ac00000000", // nolint
				Hash:     "8974b1349bc83bd18be10697a70a71199162ebdfb62993c4377ce7fdd01e025a",
				Size:     226,
				Vsize:    226,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "5385158b8721c90f47d1dd1e05aec4076f1b614c5b5b64950b8c556c8e4b6213",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3045022100ccdfdb7f9b45e5bf1d36b4963fa05714b0cb963e3f93d370f806166203904cd9022031d5d036bd92716e742756511a4be0ca85674f1d65cf8285eccb20bec9c09873[ALL] 03d3e092b233a629befb6a2249024faebfddf944140f5403a7eedb3f011c47271d", // nolint
							Hex: "483045022100ccdfdb7f9b45e5bf1d36b4963fa05714b0cb963e3f93d370f806166203904cd9022031d5d036bd92716e742756511a4be0ca85674f1d65cf8285eccb20bec9c09873012103d3e092b233a629befb6a2249024faebfddf944140f5403a7eedb3f011c47271d", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 153095.35877132,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 c5207f84f009632513b304af1e69757f79d1df45 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914c5207f84f009632513b304af1e69757f79d1df4588ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DP7Qb4KsJFcvsi4eDioLh2jMmMRcSfGJmo",
							},
						},
					},
					{
						Value: 8028.82358130,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 3d9e306694eb9a55e00113a6f9b838392d15c71a OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9143d9e306694eb9a55e00113a6f9b838392d15c71a88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DAkuG8SLbFuum36csAx9cAGxWzkRpMC3NJ",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000001e194dd90db9b7fcf661327515006bd219636ae5b6dc23120146b52ffe508753f000000006c4930460221008e6cd035285d9434b47cca0a01d6ebb4898d1a04bf16a7ef53cdda2358283edd022100b95700890afea19b08a5b67a6bc13e85df5cce990b749e23f685986b4685076e01210273c7e637ceaa0f2fd102278d48c3d4cc60b4ba2ae24a3f5daed840e515757ed6ffffffff02a1f17efc000000001976a914110e90aa49171cc316e4d7876fed69fea686491888acc0f447ac060000001976a914c62cbb894e9889bf2898e902c57ecb79b2b0254d88ac00000000", // nolint
				Hash:     "a00fdf72de70c1bc1b4a8cdb2ccf6e544aafd54e8b6fc663590c8fd410fd7aa3",
				Size:     227,
				Vsize:    227,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "3f7508e5ff526b142031c26d5bae369621bd065051271366cf7f9bdb90dd94e1",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "30460221008e6cd035285d9434b47cca0a01d6ebb4898d1a04bf16a7ef53cdda2358283edd022100b95700890afea19b08a5b67a6bc13e85df5cce990b749e23f685986b4685076e[ALL] 0273c7e637ceaa0f2fd102278d48c3d4cc60b4ba2ae24a3f5daed840e515757ed6", // nolint
							Hex: "4930460221008e6cd035285d9434b47cca0a01d6ebb4898d1a04bf16a7ef53cdda2358283edd022100b95700890afea19b08a5b67a6bc13e85df5cce990b749e23f685986b4685076e01210273c7e637ceaa0f2fd102278d48c3d4cc60b4ba2ae24a3f5daed840e515757ed6", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 42.36177825,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 110e90aa49171cc316e4d7876fed69fea6864918 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914110e90aa49171cc316e4d7876fed69fea686491888ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D6hHUX8ebMFKzNTcRZ1fDPGgN28rn3suf4",
							},
						},
					},
					{
						Value: 286.60200640,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 c62cbb894e9889bf2898e902c57ecb79b2b0254d OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914c62cbb894e9889bf2898e902c57ecb79b2b0254d88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DPCwvKABK4jzxN3QbrqDKMzMC6suJBqgqb",
							},
						},
					},
				},
			},
			{
				Hex:      "0100000001a37afd10d48f0c5963c66f8b4ed5af4a546ecf2cdb8c4a1bbcc170de72df0fa0000000006b483045022100ceb7d06c357b3cf2b83f2bd50d0e498312c85e17344497193f77e020a3874afb0220296b0e6d0609c8469e7c7f8c13e3ef2212bf61ce04e898adb70348d45819eaa30121037b2e9c9fba6e830a2a12713f89ab42e87ca34ddd9a68b06768577b4f8d543801ffffffff026a9ab816000000001976a91465c5ba865bf8b0b533fb0be75856f51db71e46cf88ac3776d0df000000001976a91420717bf29cae8d800835d4a6621a8bcd92f24f2288ac00000000", // nolint
				Hash:     "dfe755c256ac166db49654a6bf76d637dab24857bdff775b49b26a5a7a6fbd3b",
				Size:     226,
				Vsize:    226,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "a00fdf72de70c1bc1b4a8cdb2ccf6e544aafd54e8b6fc663590c8fd410fd7aa3",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3045022100ceb7d06c357b3cf2b83f2bd50d0e498312c85e17344497193f77e020a3874afb0220296b0e6d0609c8469e7c7f8c13e3ef2212bf61ce04e898adb70348d45819eaa3[ALL] 037b2e9c9fba6e830a2a12713f89ab42e87ca34ddd9a68b06768577b4f8d543801", // nolint
							Hex: "483045022100ceb7d06c357b3cf2b83f2bd50d0e498312c85e17344497193f77e020a3874afb0220296b0e6d0609c8469e7c7f8c13e3ef2212bf61ce04e898adb70348d45819eaa30121037b2e9c9fba6e830a2a12713f89ab42e87ca34ddd9a68b06768577b4f8d543801", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 3.81196906,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 65c5ba865bf8b0b533fb0be75856f51db71e46cf OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a91465c5ba865bf8b0b533fb0be75856f51db71e46cf88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DERDf4KwNAepVagbNbQ1LRMe6n5bjreMMr",
							},
						},
					},
					{
						Value: 37.54980919,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 20717bf29cae8d800835d4a6621a8bcd92f24f22 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a91420717bf29cae8d800835d4a6621a8bcd92f24f2288ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"D86e7R3kdLhk849oC3kzbbVQw8yt5WWw49",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000015a021ed0fde77c37c49329b6dfeb629119710aa79706e18bd13bc89b34b17489000000006c493046022100f766c02b8be767b5e40ab2e5fca085b095883b24e76da75a7894f2134f85e2f60221008bb3fa3322b9f032eab5fbb9c3c3792629c5444a628e489bbef074507d2822bd0121039591af3b189bc9063110783ec97b1e061cd88dea2b374efcecdb2b4fdaf48d1affffffff02d8baa4478f0d00001976a9149fdbc530d583f6e0dc31e9281db0c9db39325e5a88ac34fcd7395d0000001976a914770d36fdf8a25cb438fed008540ae268e62addbf88ac00000000", // nolint
				Hash:     "b3a66d10f35a4de16333386ac343b99d534ba26633f725a8cf5fc40d120d60d5",
				Size:     227,
				Vsize:    227,
				Version:  1,
				Locktime: 0,
				Weight:   0,
				Inputs: []*Input{
					{
						TxHash: "8974b1349bc83bd18be10697a70a71199162ebdfb62993c4377ce7fdd01e025a",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3046022100f766c02b8be767b5e40ab2e5fca085b095883b24e76da75a7894f2134f85e2f60221008bb3fa3322b9f032eab5fbb9c3c3792629c5444a628e489bbef074507d2822bd[ALL] 039591af3b189bc9063110783ec97b1e061cd88dea2b374efcecdb2b4fdaf48d1a", // nolint
							Hex: "493046022100f766c02b8be767b5e40ab2e5fca085b095883b24e76da75a7894f2134f85e2f60221008bb3fa3322b9f032eab5fbb9c3c3792629c5444a628e489bbef074507d2822bd0121039591af3b189bc9063110783ec97b1e061cd88dea2b374efcecdb2b4fdaf48d1a", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 149090.33462488,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 9fdbc530d583f6e0dc31e9281db0c9db39325e5a OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a9149fdbc530d583f6e0dc31e9281db0c9db39325e5a88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DKiMFfacJsqzW4pF2LvtwNdQKiu22XwrxR",
							},
						},
					},
					{
						Value: 4004.02414644,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 770d36fdf8a25cb438fed008540ae268e62addbf OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914770d36fdf8a25cb438fed008540ae268e62addbf88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"DFzanGUHhPynMq3ompvBRjqyvBS5mFFStg",
							},
						},
					},
				},
			},
			{
				Hash:     "fake",
				Hex:      "fake hex",
				Version:  2,
				Size:     421,
				Vsize:    612,
				Weight:   129992,
				Locktime: 10,
				Inputs: []*Input{
					{
						TxHash: "503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "00142b2296c588ec413cebd19c3cbc04ea830ead6e78",
							Hex: "1600142b2296c588ec413cebd19c3cbc04ea830ead6e78",
						},
						TxInWitness: []string{
							"304402205f39ccbab38b644acea0776d18cb63ce3e37428cbac06dc23b59c61607aef69102206b8610827e9cb853ea0ba38983662034bd3575cc1ab118fb66d6a98066fa0bed01", // nolint
							"0304c01563d46e38264283b99bb352b46e69bf132431f102d4bd9a9d8dab075e7f",
						},
						Sequence: 4294967295,
					},
					{
						TxHash: "503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "00142b2296c588ec413cebd19c3cbc04ea830ead6e78",
							Hex: "1600142b2296c588ec413cebd19c3cbc04ea830ead6e78",
						},
						Sequence: 4294967295,
					},
					{
						TxHash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
							Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 200.56,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
							},
						},
					},
				},
			},
		},
	}
)

func TestNetworkStatus(t *testing.T) {
	tests := map[string]struct {
		responses []responseFixture

		expectedStatus *types.NetworkStatusResponse
		expectedError  error
	}{
		"successful": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_blockchain_info_response.json"),
					url:    url,
				},
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_response.json"),
					url:    url,
				},
				{
					status: http.StatusOK,
					body:   loadFixture("get_peer_info_response.json"),
					url:    url,
				},
			},
			expectedStatus: &types.NetworkStatusResponse{
				CurrentBlockIdentifier: blockIdentifier1000,
				CurrentBlockTimestamp:  block1000.Time * 1000,
				GenesisBlockIdentifier: MainnetGenesisBlockIdentifier,
				Peers: []*types.Peer{
					{
						PeerID: "34.221.250.46:22556",
						Metadata: forceMarshalMap(t, &PeerInfo{
							Addr:           "34.221.250.46:22556",
							Version:        70015,
							SubVer:         "/Shibetoshi:1.14.2/",
							StartingHeight: 3667937,
							RelayTxes:      true,
							LastSend:       1617139759,
							LastRecv:       1617139785,
							BanScore:       0,
							SyncedHeaders:  -1,
							SyncedBlocks:   -1,
						}),
					},
					{
						PeerID: "54.38.205.113:22556",
						Metadata: forceMarshalMap(t, &PeerInfo{
							Addr:           "54.38.205.113:22556",
							RelayTxes:      true,
							LastSend:       1617139795,
							LastRecv:       1617139795,
							Version:        70015,
							SubVer:         "/Shibetoshi:1.14.2/",
							StartingHeight: 3667937,
							BanScore:       0,
							SyncedHeaders:  -1,
							SyncedBlocks:   -1,
						}),
					},
				},
			},
		},
		"blockchain warming up error": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("rpc_in_warmup_response.json"),
					url:    url,
				},
			},
			expectedError: errors.New("rpc in warmup"),
		},
		"blockchain info error": {
			responses: []responseFixture{
				{
					status: http.StatusInternalServerError,
					body:   "{}",
					url:    url,
				},
			},
			expectedError: errors.New("invalid response: 500 Internal Server Error"),
		},
		"peer info not accessible": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_blockchain_info_response.json"),
					url:    url,
				},
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_response.json"),
					url:    url,
				},
				{
					status: http.StatusInternalServerError,
					body:   "{}",
					url:    url,
				},
			},
			expectedError: errors.New("invalid response: 500 Internal Server Error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				assert = assert.New(t)
			)

			responses := make(chan responseFixture, len(test.responses))
			for _, response := range test.responses {
				responses <- response
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := <-responses
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("POST", r.Method)
				assert.Equal(response.url, r.URL.RequestURI())

				w.WriteHeader(response.status)
				fmt.Fprintln(w, response.body)
			}))

			client := NewClient(ts.URL, MainnetGenesisBlockIdentifier, MainnetCurrency)
			status, err := client.NetworkStatus(context.Background())
			if test.expectedError != nil {
				assert.Contains(err.Error(), test.expectedError.Error())
			} else {
				assert.NoError(err)
				assert.Equal(test.expectedStatus, status)
			}
		})
	}
}

func TestGetPeers(t *testing.T) {
	tests := map[string]struct {
		responses []responseFixture

		expectedPeers []*types.Peer
		expectedError error
	}{
		"successful": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_peer_info_response.json"),
					url:    url,
				},
			},
			expectedPeers: []*types.Peer{
				{
					PeerID: "34.221.250.46:22556",
					Metadata: forceMarshalMap(t, &PeerInfo{
						Addr:           "34.221.250.46:22556",
						Version:        70015,
						SubVer:         "/Shibetoshi:1.14.2/",
						StartingHeight: 3667937,
						RelayTxes:      true,
						LastSend:       1617139759,
						LastRecv:       1617139785,
						BanScore:       0,
						SyncedHeaders:  -1,
						SyncedBlocks:   -1,
					}),
				},
				{
					PeerID: "54.38.205.113:22556",
					Metadata: forceMarshalMap(t, &PeerInfo{
						Addr:           "54.38.205.113:22556",
						RelayTxes:      true,
						LastSend:       1617139795,
						LastRecv:       1617139795,
						Version:        70015,
						SubVer:         "/Shibetoshi:1.14.2/",
						StartingHeight: 3667937,
						BanScore:       0,
						SyncedHeaders:  -1,
						SyncedBlocks:   -1,
					}),
				},
			},
		},
		"blockchain warming up error": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("rpc_in_warmup_response.json"),
					url:    url,
				},
			},
			expectedError: errors.New("rpc in warmup"),
		},
		"peer info error": {
			responses: []responseFixture{
				{
					status: http.StatusInternalServerError,
					body:   "{}",
					url:    url,
				},
			},
			expectedError: errors.New("invalid response: 500 Internal Server Error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				assert = assert.New(t)
			)

			responses := make(chan responseFixture, len(test.responses))
			for _, response := range test.responses {
				responses <- response
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := <-responses
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("POST", r.Method)
				assert.Equal(response.url, r.URL.RequestURI())

				w.WriteHeader(response.status)
				fmt.Fprintln(w, response.body)
			}))

			client := NewClient(ts.URL, MainnetGenesisBlockIdentifier, MainnetCurrency)
			peers, err := client.GetPeers(context.Background())
			if test.expectedError != nil {
				assert.Contains(err.Error(), test.expectedError.Error())
			} else {
				assert.NoError(err)
				assert.Equal(test.expectedPeers, peers)
			}
		})
	}
}

func TestGetRawBlock(t *testing.T) {
	tests := map[string]struct {
		blockIdentifier *types.PartialBlockIdentifier
		responses       []responseFixture

		expectedBlock *Block
		expectedCoins []string
		expectedError error
	}{
		"lookup by hash": {
			blockIdentifier: &types.PartialBlockIdentifier{
				Hash: &blockIdentifier1000.Hash,
			},
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_response.json"),
					url:    url,
				},
			},
			expectedBlock: block1000,
			expectedCoins: []string{},
		},
		"lookup by hash 2": {
			blockIdentifier: &types.PartialBlockIdentifier{
				Hash: &blockIdentifier100000.Hash,
			},
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_response_2.json"),
					url:    url,
				},
			},
			expectedBlock: block100000,
			expectedCoins: []string{
				"87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
				"503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5:0",
				"503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5:1",
			},
		},
		"lookup by hash (get block api error)": {
			blockIdentifier: &types.PartialBlockIdentifier{
				Hash: &blockIdentifier1000.Hash,
			},
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_not_found_response.json"),
					url:    url,
				},
			},
			expectedError: ErrBlockNotFound,
		},
		"lookup by hash (get block internal error)": {
			blockIdentifier: &types.PartialBlockIdentifier{
				Hash: &blockIdentifier1000.Hash,
			},
			responses: []responseFixture{
				{
					status: http.StatusInternalServerError,
					body:   "{}",
					url:    url,
				},
			},
			expectedBlock: nil,
			expectedError: errors.New("invalid response: 500 Internal Server Error"),
		},
		"lookup by index": {
			blockIdentifier: &types.PartialBlockIdentifier{
				Index: &blockIdentifier1000.Index,
			},
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_hash_response.json"),
					url:    url,
				},
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_response.json"),
					url:    url,
				},
			},
			expectedBlock: block1000,
			expectedCoins: []string{},
		},
		"lookup by index (out of range)": {
			blockIdentifier: &types.PartialBlockIdentifier{
				Index: &blockIdentifier1000.Index,
			},
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_hash_out_of_range_response.json"),
					url:    url,
				},
			},
			expectedError: errors.New("height out of range"),
		},
		"current block lookup": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_blockchain_info_response.json"),
					url:    url,
				},
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_response.json"),
					url:    url,
				},
			},
			expectedBlock: block1000,
			expectedCoins: []string{},
		},
		"current block lookup (can't get current info)": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("rpc_in_warmup_response.json"),
					url:    url,
				},
			},
			expectedError: errors.New("unable to get blockchain info"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				assert = assert.New(t)
			)

			responses := make(chan responseFixture, len(test.responses))
			for _, response := range test.responses {
				responses <- response
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := <-responses
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("POST", r.Method)
				assert.Equal(response.url, r.URL.RequestURI())

				w.WriteHeader(response.status)
				fmt.Fprintln(w, response.body)
			}))

			client := NewClient(ts.URL, MainnetGenesisBlockIdentifier, MainnetCurrency)
			block, coins, err := client.GetRawBlock(context.Background(), test.blockIdentifier)
			if test.expectedError != nil {
				assert.Contains(err.Error(), test.expectedError.Error())
			} else {
				assert.NoError(err)
				assert.Equal(test.expectedBlock, block)
				assert.Equal(test.expectedCoins, coins)
			}
		})
	}
}

func int64Pointer(v int64) *int64 {
	return &v
}

func mustMarshalMap(v interface{}) map[string]interface{} {
	m, _ := types.MarshalMap(v)
	return m
}

func TestParseBlock(t *testing.T) {
	tests := map[string]struct {
		block *Block
		coins map[string]*types.AccountCoin

		expectedBlock *types.Block
		expectedError error
	}{
		"no fetched transactions": {
			block: block1000,
			coins: map[string]*types.AccountCoin{},
			expectedBlock: &types.Block{
				BlockIdentifier: blockIdentifier1000,
				ParentBlockIdentifier: &types.BlockIdentifier{
					Hash:  "a375b78d24b8cd5d554da5d7ddc787412dbe963ea6ba6802aab46452a475066e",
					Index: 999,
				},
				Timestamp: 1386481090,
				Transactions: []*types.Transaction{
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "9480ac41aac2674ac498849b9ab95661ee73fb372140e62c6e7a6fa29f5a09d1",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   CoinbaseOpType,
								Status: types.String(SuccessStatus),
								Metadata: mustMarshalMap(&OperationMetadata{
									Coinbase: "04ca05a4520131062f503253482f",
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "DG3FAujuBTvozAU423arUKYnhQQa6XzWZ9", // nolint
								},
								Amount: &types.Amount{
									Value:    "185878.00000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "9480ac41aac2674ac498849b9ab95661ee73fb372140e62c6e7a6fa29f5a09d1:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "03236a82a90eda514a3373f1a13b349cbff983d00879023ce15ef5ff8c757aa7db OP_CHECKSIG", // nolint
										Hex:  "2103236a82a90eda514a3373f1a13b349cbff983d00879023ce15ef5ff8c757aa7dbac",         // nolint
										Type: "pubkey",
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    109,
							Version: 1,
							Vsize:   109,
							Weight:  0,
						}),
					},
				},
				Metadata: mustMarshalMap(&BlockMetadata{
					Size:       190,
					Weight:     760,
					Version:    1,
					MerkleRoot: "9480ac41aac2674ac498849b9ab95661ee73fb372140e62c6e7a6fa29f5a09d1",
					MedianTime: 1386480992,
					Nonce:      2308638208,
					Bits:       "1d03d07b",
					Difficulty: 0.2621620216098152,
				}),
			},
		},
		"block 100000": {
			block: block100000,
			coins: map[string]*types.AccountCoin{
				"c2e410a0c9ff9ee2808a5efc27885e092baf337290dca77faf1b13da5a946d98:0": {
					Account: &types.AccountIdentifier{
						Address: "DCvXWEMbJNFMuZjDcBTVd8nG6K2kCDTKQu",
					},
					Coin: &types.Coin{
						CoinIdentifier: &types.CoinIdentifier{
							Identifier: "c2e410a0c9ff9ee2808a5efc27885e092baf337290dca77faf1b13da5a946d98:0",
						},
						Amount: &types.Amount{
							Value:    "42233775.22310657",
							Currency: MainnetCurrency,
						},
					},
				},
				"503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5:0": {
					Account: &types.AccountIdentifier{
						Address: "3FfQGY7jqsADC7uTVqF3vKQzeNPiBPTqt4",
					},
					Coin: &types.Coin{
						CoinIdentifier: &types.CoinIdentifier{
							Identifier: "503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5:0",
						},
						Amount: &types.Amount{
							Value:    "3467607",
							Currency: MainnetCurrency,
						},
					},
				},
				"503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5:1": {
					Account: &types.AccountIdentifier{
						Address: "1NdvAyRJLdK5EXs7DV3ebYb5wffdCZk1pD",
					},
					Coin: &types.Coin{
						CoinIdentifier: &types.CoinIdentifier{
							Identifier: "503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5:1",
						},
						Amount: &types.Amount{
							Value:    "0",
							Currency: MainnetCurrency,
						},
					},
				},
			},
			expectedBlock: &types.Block{
				BlockIdentifier: blockIdentifier100000,
				ParentBlockIdentifier: &types.BlockIdentifier{
					Hash:  "12aca0938fe1fb786c9e0e4375900e8333123de75e240abd3337d1b411d14ebe",
					Index: 99999,
				},
				Timestamp: 1392346405,
				Transactions: []*types.Transaction{
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "8c14f0db3df150123e6f3dbbf30f8b955a8249b62ac1d1ff16284aefa3d06d87",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   CoinbaseOpType,
								Status: types.String(SuccessStatus),
								Metadata: mustMarshalMap(&OperationMetadata{
									Coinbase: "044c86041b020602",
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "34qkc2iac6RsyxZVfyE2S5U5WcRsbg2dpK",
								},
								Amount: &types.Amount{
									Value:    "1589351625",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "8c14f0db3df150123e6f3dbbf30f8b955a8249b62ac1d1ff16284aefa3d06d87:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_HASH160 228f554bbf766d6f9cc828de1126e3d35d15e5fe OP_EQUAL",
										Hex:          "a914228f554bbf766d6f9cc828de1126e3d35d15e5fe87",
										RequiredSigs: 1,
										Type:         "scripthash",
										Addresses: []string{
											"34qkc2iac6RsyxZVfyE2S5U5WcRsbg2dpK",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "6a24aa21a9ed10109f4b82aa3ed7ec9d02a2a90246478b3308c8b85daf62fe501d58d05727a4",
								},
								Amount: &types.Amount{
									Value:    "0",
									Currency: MainnetCurrency,
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_RETURN aa21a9ed10109f4b82aa3ed7ec9d02a2a90246478b3308c8b85daf62fe501d58d05727a4",
										Hex:  "6a24aa21a9ed10109f4b82aa3ed7ec9d02a2a90246478b3308c8b85daf62fe501d58d05727a4",
										Type: "nulldata",
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    135,
							Version: 1,
							Vsize:   135,
							Weight:  540,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fake",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-3467607",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "3FfQGY7jqsADC7uTVqF3vKQzeNPiBPTqt4",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "00142b2296c588ec413cebd19c3cbc04ea830ead6e78",
										Hex: "1600142b2296c588ec413cebd19c3cbc04ea830ead6e78",
									},
									TxInWitness: []string{
										"304402205f39ccbab38b644acea0776d18cb63ce3e37428cbac06dc23b59c61607aef69102206b8610827e9cb853ea0ba38983662034bd3575cc1ab118fb66d6a98066fa0bed01", // nolint
										"0304c01563d46e38264283b99bb352b46e69bf132431f102d4bd9a9d8dab075e7f",
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(1),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "0",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1NdvAyRJLdK5EXs7DV3ebYb5wffdCZk1pD",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "00142b2296c588ec413cebd19c3cbc04ea830ead6e78",
										Hex: "1600142b2296c588ec413cebd19c3cbc04ea830ead6e78",
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(2),
								},
								Type:   InputOpType,
								Status: types.String(SuccessStatus),
								Amount: &types.Amount{
									Value:    "-556000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        3,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: types.String(SuccessStatus),
								Account: &types.AccountIdentifier{
									Address: "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
								},
								Amount: &types.Amount{
									Value:    "20056000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fake:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:     421,
							Version:  2,
							Vsize:    612,
							Weight:   129992,
							Locktime: 10,
						}),
					},
				},
				Metadata: mustMarshalMap(&BlockMetadata{
					Size:       13372,
					Weight:     53488,
					Version:    1,
					MerkleRoot: "31757c266102d1bee62ef2ff8438663107d64bdd5d9d9173421ec25fb2a814de",
					MedianTime: 1392346434,
					Nonce:      2216773632,
					Bits:       "1b267eeb",
					Difficulty: 1702.39468793143,
				}),
			},
		},
		"missing transactions": {
			block:         block100000,
			coins:         map[string]*types.AccountCoin{},
			expectedError: errors.New("error finding previous tx"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				assert = assert.New(t)
			)

			client := NewClient("", MainnetGenesisBlockIdentifier, MainnetCurrency)
			block, err := client.ParseBlock(context.Background(), test.block, test.coins)
			if test.expectedError != nil {
				assert.Contains(err.Error(), test.expectedError.Error())
			} else {
				assert.NoError(err)
				assert.Equal(test.expectedBlock, block)
			}
		})
	}
}

func TestSuggestedFeeRate(t *testing.T) {
	tests := map[string]struct {
		responses []responseFixture

		expectedRate  float64
		expectedError error
	}{
		"successful": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("fee_rate.json"),
					url:    url,
				},
			},
			expectedRate: float64(0.00001),
		},
		"invalid range error": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("invalid_fee_rate.json"),
					url:    url,
				},
			},
			expectedError: errors.New("error getting fee estimate"),
		},
		"500 error": {
			responses: []responseFixture{
				{
					status: http.StatusInternalServerError,
					body:   "{}",
					url:    url,
				},
			},
			expectedError: errors.New("invalid response: 500 Internal Server Error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				assert = assert.New(t)
			)

			responses := make(chan responseFixture, len(test.responses))
			for _, response := range test.responses {
				responses <- response
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := <-responses
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("POST", r.Method)
				assert.Equal(response.url, r.URL.RequestURI())

				w.WriteHeader(response.status)
				fmt.Fprintln(w, response.body)
			}))

			client := NewClient(ts.URL, MainnetGenesisBlockIdentifier, MainnetCurrency)
			rate, err := client.SuggestedFeeRate(context.Background(), 1)
			if test.expectedError != nil {
				assert.Contains(err.Error(), test.expectedError.Error())
			} else {
				assert.NoError(err)
				assert.Equal(test.expectedRate, rate)
			}
		})
	}
}

func TestRawMempool(t *testing.T) {
	tests := map[string]struct {
		responses []responseFixture

		expectedTransactions []string
		expectedError        error
	}{
		"successful": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("raw_mempool.json"),
					url:    url,
				},
			},
			expectedTransactions: []string{
				"9cec12d170e97e21a876fa2789e6bfc25aa22b8a5e05f3f276650844da0c33ab",
				"37b4fcc8e0b229412faeab8baad45d3eb8e4eec41840d6ac2103987163459e75",
				"7bbb29ae32117597fcdf21b464441abd571dad52d053b9c2f7204f8ea8c4762e",
			},
		},
		"500 error": {
			responses: []responseFixture{
				{
					status: http.StatusInternalServerError,
					body:   "{}",
					url:    url,
				},
			},
			expectedError: errors.New("invalid response: 500 Internal Server Error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				assert = assert.New(t)
			)

			responses := make(chan responseFixture, len(test.responses))
			for _, response := range test.responses {
				responses <- response
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := <-responses
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("POST", r.Method)
				assert.Equal(response.url, r.URL.RequestURI())

				w.WriteHeader(response.status)
				fmt.Fprintln(w, response.body)
			}))

			client := NewClient(ts.URL, MainnetGenesisBlockIdentifier, MainnetCurrency)
			txs, err := client.RawMempool(context.Background())
			if test.expectedError != nil {
				assert.Contains(err.Error(), test.expectedError.Error())
			} else {
				assert.NoError(err)
				assert.Equal(test.expectedTransactions, txs)
			}
		})
	}
}

// loadFixture takes a file name and returns the response fixture.
func loadFixture(fileName string) string {
	content, err := ioutil.ReadFile(fmt.Sprintf("client_fixtures/%s", fileName))
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

type responseFixture struct {
	status int
	body   string
	url    string
}
