/*******************************************************************************
 * Copyright 2020 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package test

import (
	"fmt"
	"strings"
)

var privateKey = fmt.Sprintf("%s",
	"-----BEGIN RSA PRIVATE KEY-----\n"+
		"MIIEpQIBAAKCAQEA+Cl+aPrXYrmBLPj7MflVSnC0Cnb+08C8+UyQx/llWbdzbsOp\n"+
		"pUix2gvMeUIeqbrDggu1ac8ogwjRuRdiAL0sygCutCMD4jB/2d0x2FEgvKbRdufD\n"+
		"2tbOQ2PniqoJr9IbYHMJYWAKDKWrCXcV+g5Jk/hKSCtmDrgvGN3gmbargtWmFGLb\n"+
		"IGcJ9HqfAeLc4tuCZ91i/I1cvKUaqfdAJkJeEbrHw0PuGDywpqfPlDCclJX+cOOP\n"+
		"W4jCGWd4Db42vlSjCBfq4Z0MBX++u4P8uyx6SC2LrlEdRO43yB/ryBR17IioCSxu\n"+
		"khavfZHbYPEwrLTYarx8FY62FUgcffzcsmhxZwIDAQABAoIBAQCtw1gpJ+Mi1KOX\n"+
		"mutAzcYj7pCSd0nteZqYsTz7WSzXSjYAi+7Atgsak4JkMaEI1aZJ6+rmINDMF6PK\n"+
		"B45u2AeBlkK+DXqNqcoMAe8B+aSDlAc9TAF+vUQGOfEJzhAkVWkn+sTJsxa2TlZZ\n"+
		"tVHlGpX4jzVsHT9D9UG9FrdKynaDj7bMf+sJHiuo9cAAIauZ3BTIAR81Zu+SgfYI\n"+
		"95MqW9qrvBVW0zIlqkNeYXbM6dCGZlVm2BS2EdcfwFvhbKjXFqUUm8FOn4VQtcSs\n"+
		"5/F4oDSncb6uT4hyC2335A7EvEKcfkaSb1AbpDUMIhjvlWfYeyXlQvYUIpKCzE21\n"+
		"WnkILCCJAoGBAP+nnPDMq9nklp3ZBWMzF4wvdZD3nxep4oFg7fwIPKNnEHpI3/zI\n"+
		"fDQsjnMpvzP/ZOzmMXr0sE+XwiDJ4nCCRTm0+yCD3P1VZ4KKbp6MaWa5pRCCk9QK\n"+
		"FpnzZu4C1l4JgG+HJilxN/jSfijRsxLcmbRm2fEFbjeQdV3sj9UHNHo1AoGBAPh/\n"+
		"SlKCOnSUJGCLZzbj3dUG1uso0UzE9mY8NoZxS5L1vSh98OejjXJb9VmAafsZmjdx\n"+
		"PMd/sggLkgeQTkIu4AUEwsX9i/9bOzjyXr3+SPxF+BIc5LH4zuFXoyiCiFUaYmfB\n"+
		"XvwXs9OO7/F//oRiFVY1uFavKnWvi7gr1k8CyZCrAoGBAJaDX+qFFUgbRHF6K6nT\n"+
		"krF934GRx6Bu7GOvZW1UjB7HtvPHo9d3UWiGMveqRF+gpRK0E72IAaVae3hCY4ZJ\n"+
		"q+flnVPvTlP3zBEW3zmJASTxdzTZK59SsSvCGX9XPE3w2iTPNLCBb6qWgqAVlZAt\n"+
		"QHDtfLJhuBoOeorpk2Sf8U1hAoGBAO5MKvqqleH7ulK2/DjAFZfWoj0KfIPREbUC\n"+
		"owsUFHQOoeH1vBJ2Xgs/si2tHnTEnYXzWmS5yQE8D0KfmNyQ1RUa9qklNp6fX1CB\n"+
		"5GbwNg9uDbFY8drVjZa9EuKjIpfx4FI9NpgrJrCHDwQZSPqskGeGxoqiGeaXfDYW\n"+
		"G8LTGnZXAoGAScksu0S+ScRp+++9KvlyTLD9uMw3Soceb/iqA63ZZI814tlQO10K\n"+
		"/0kXqd4O+fR5Ef7+J6f68RQ8ND4XmJJy8lcE4vkc1DS3PYznyLB/2tPZzeNxYOCE\n"+
		"8Pi5JgSArgEYUwWJtob3rTsUoBjnEB76AtgIXCDBjo/GWorhvTyPMJw=\n"+
		"-----END RSA PRIVATE KEY-----\n",
)

var ValidPrivateKey = []byte(privateKey)

var InvalidPrivateKey = []byte(strings.ToUpper(privateKey))

var publicKey = fmt.Sprintf("%s",
	"-----BEGIN PUBLIC KEY-----\n"+
		"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+Cl+aPrXYrmBLPj7MflV\n"+
		"SnC0Cnb+08C8+UyQx/llWbdzbsOppUix2gvMeUIeqbrDggu1ac8ogwjRuRdiAL0s\n"+
		"ygCutCMD4jB/2d0x2FEgvKbRdufD2tbOQ2PniqoJr9IbYHMJYWAKDKWrCXcV+g5J\n"+
		"k/hKSCtmDrgvGN3gmbargtWmFGLbIGcJ9HqfAeLc4tuCZ91i/I1cvKUaqfdAJkJe\n"+
		"EbrHw0PuGDywpqfPlDCclJX+cOOPW4jCGWd4Db42vlSjCBfq4Z0MBX++u4P8uyx6\n"+
		"SC2LrlEdRO43yB/ryBR17IioCSxukhavfZHbYPEwrLTYarx8FY62FUgcffzcsmhx\n"+
		"ZwIDAQAB\n"+
		"-----END PUBLIC KEY-----\n",
)

var ValidPublicKey = []byte(publicKey)

var InvalidPublicKey = []byte(strings.ToUpper(publicKey))
