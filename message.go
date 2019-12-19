package netfilter

import (
	"github.com/pkg/errors"

	"github.com/mdlayher/netlink"
)

// UnmarshalNetlink unmarshals a netlink.Message into a Netfilter Header and Attributes.
func UnmarshalNetlink(msg netlink.Message) (Header, []Attribute, error) {

	h, ad, err := DecodeNetlink(msg)
	if err != nil {
		return Header{}, nil, err
	}

	attrs, err := unmarshalAttributes(ad)
	if err != nil {
		return Header{}, nil, err
	}

	return h, attrs, nil
}

// DecodeNetlink returns msg's Netfilter header and an AttributeDecoder that can be used
// to iteratively decode all Netlink attributes contained in the message.
func DecodeNetlink(msg netlink.Message) (Header, *netlink.AttributeDecoder, error) {

	var h Header

	err := h.unmarshal(msg)
	if err != nil {
		return Header{}, nil, errors.Wrap(err, "unmarshaling netfilter header")
	}

	ad, err := NewAttributeDecoder(msg.Data[nfHeaderLen:])
	if err != nil {
		return Header{}, nil, errors.Wrap(err, "creating attribute decoder")
	}

	return h, ad, nil
}

// MarshalNetlink takes a Netfilter Header and Attributes and returns a netlink.Message.
func MarshalNetlink(h Header, attrs []Attribute) (netlink.Message, error) {

	ba, err := MarshalAttributes(attrs)
	if err != nil {
		return netlink.Message{}, err
	}

	// initialize with 4 bytes of Data before unmarshal
	nlm := netlink.Message{Data: make([]byte, 4)}

	// marshal error ignored, safe to do if msg Data is initialized
	_ = h.marshal(&nlm)

	nlm.Data = append(nlm.Data, ba...)

	return nlm, nil
}
