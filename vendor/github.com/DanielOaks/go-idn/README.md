# Go-idn

Go-idn is a project that hopes to bring IDN and Stringprep to Go and aims to become feature compatible with libidn.

This library is in a VERY EARLY stage, since I've only just forked and updated it. Things WILL CHANGE and things MAY NOT WORK properly yet.

---

[![Build Status](https://travis-ci.org/DanielOaks/go-idn.svg?branch=master)](https://travis-ci.org/DanielOaks/go-idn) [![Coverage Status](https://coveralls.io/repos/github/DanielOaks/go-idn/badge.svg?branch=master)](https://coveralls.io/github/DanielOaks/go-idn?branch=master)

---

Go-idn is a mostly-documented implementation of the Stringprep, Punycode and IDNA specifications. Go-idn's purpose is to encode and decode internationalized domain names and provide a simple Stringprep interface using pure Go code.

The library contains a generic Stringprep implementation. Profiles for Nameprep are included, and we plan to support iSCSI, SASL and XMPP profiles. Punycode and ASCII Compatible Encoding (ACE) via IDNA are supported. A mechanism to define Top-Level Domain (TLD) specific validation tables, and to compare strings against those tables, is included. Default tables for some TLDs are also included. 
