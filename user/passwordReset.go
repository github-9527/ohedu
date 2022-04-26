package user

import (
	"fmt"

	ldap "github.com/go-ldap/ldap"
	"golang.org/x/text/encoding/unicode"
	ber "gopkg.in/asn1-ber.v1"
)

const (
	ldapAttrAccountName                        = "sAMAccountName"
	ldapAttrDN                                 = "dn"
	ldapAttrUAC                                = "userAccountControl"
	ldapAttrUPN                                = "userPrincipalName" // username@logon.domain
	ldapAttrEmail                              = "mail"
	ldapAttrUnicodePw                          = "unicodePwd"
	controlTypeLdapServerPolicyHints           = "1.2.840.113556.1.4.2239"
	controlTypeLdapServerPolicyHintsDeprecated = "1.2.840.113556.1.4.2066"
)

type (
	// ldapControlServerPolicyHints implements ldap.Control
	ldapControlServerPolicyHints struct {
		oid string
	}
)

// GetControlType implements ldap.Control
func (c *ldapControlServerPolicyHints) GetControlType() string {
	return c.oid
}

// Encode implements ldap.Control
func (c *ldapControlServerPolicyHints) Encode() *ber.Packet {
	packet := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "Control")
	packet.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, c.GetControlType(), "Control Type (LDAP_SERVER_POLICY_HINTS_OID)"))
	packet.AppendChild(ber.NewBoolean(ber.ClassUniversal, ber.TypePrimitive, ber.TagBoolean, true, "Criticality"))

	p2 := ber.Encode(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, nil, "Control Value (Policy Hints)")
	seq := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "PolicyHintsRequestValue")
	seq.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagInteger, 1, "Flags"))
	p2.AppendChild(seq)
	packet.AppendChild(p2)

	return packet
}

// String implements ldap.Control
func (c *ldapControlServerPolicyHints) String() string {
	return "Enforce password history policies during password set: " + c.GetControlType()
}

// ChangePassword modifies the user password of a user
func ChangePassword(userdn string, password string, conn *ldap.Conn) error {
	//

	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	// The password needs to be enclosed in quotes
	pwdEncoded, err := utf16.NewEncoder().String(fmt.Sprintf("\"%s\"", password))
	if err != nil {
		return err
	}

	// add additional control to request if supported
	controlTypes, err := getSupportedControl(conn)
	if err != nil {
		return err
	}
	control := []ldap.Control{}
	for _, oid := range controlTypes {
		if oid == controlTypeLdapServerPolicyHints || oid == controlTypeLdapServerPolicyHintsDeprecated {
			control = append(control, &ldapControlServerPolicyHints{oid: oid})
			break
		}
	}

	passReq := ldap.NewModifyRequest(userdn, control)
	passReq.Replace(ldapAttrUnicodePw, []string{pwdEncoded})
	return conn.Modify(passReq)
}

// getSupportedControl retrieves supported extended control types
func getSupportedControl(conn ldap.Client) ([]string, error) {
	req := ldap.NewSearchRequest("", ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false, "(objectClass=*)", []string{"supportedControl"}, nil)
	res, err := conn.Search(req)
	if err != nil {
		return nil, err
	}
	return res.Entries[0].GetAttributeValues("supportedControl"), nil
}
