package zebra

import (
	"fmt"
	"io"
	"strings"

	"github.com/glycerine/zebrapack/msgp"
)

func ZkindFromString(s string) Zkind {
	s = strings.ToLower(s)
	switch s {
	case "":
		return Invalid
	case "invalid":
		return Invalid
	case "bytes":
		return Bytes
	case "string":
		return String
	case "float32":
		return Float32
	case "float64":
		return Float64
	case "complex64":
		return Complex64
	case "complex128":
		return Complex128
	case "uint":
		return Uint
	case "uint8":
		return Uint8
	case "uint16":
		return Uint16
	case "uint32":
		return Uint32
	case "uint64":
		return Uint64
	case "byte":
		return Byte
	case "int":
		return Int
	case "int8":
		return Int8
	case "int16":
		return Int16
	case "int32":
		return Int32
	case "int64":
		return Int64
	case "bool":
		return Bool
	case "intf":
		return Intf
	case "time":
		return Time
	case "ext":
		return Ext
	case "ident":
		// IDENT typically means a named struct
		return IDENT
	case "baseelem":
		return BaseElemCat
	case "map":
		return MapCat
	case "struct":
		return StructCat
	case "slice":
		return SliceCat
	case "array":
		return ArrayCat
	case "pointer":
		return PointerCat
	}
	panic(fmt.Errorf("unrecognized arg '%s' to ZkindFromString()", s))
}

func (i Zkind) String() string {
	switch i {
	case Invalid:
		return ""
	case Bytes:
		return "bytes"
	case String:
		return "string"
	case Float32:
		return "float32"
	case Float64:
		return "float64"
	case Complex64:
		return "complex64"
	case Complex128:
		return "complex128"
	case Uint:
		return "uint"
	case Uint8:
		return "uint8"
	case Uint16:
		return "uint16"
	case Uint32:
		return "uint32"
	case Uint64:
		return "uint64"
	case Byte:
		return "byte"
	case Int:
		return "int"
	case Int8:
		return "int8"
	case Int16:
		return "int16"
	case Int32:
		return "int32"
	case Int64:
		return "int64"
	case Bool:
		return "bool"

		// compound/non-primitives are uppercased
		// for readability
	case Intf:
		return "Intf"
	case Time:
		return "Time"
	case Ext:
		return "Ext"
	case IDENT:
		// IDENT typically means a named struct
		return "IDENT"
	case BaseElemCat:
		return "BaseElem"
	case MapCat:
		return "Map"
	case StructCat:
		return "Struct"
	case SliceCat:
		return "Slice"
	case ArrayCat:
		return "Array"
	case PointerCat:
		return "Pointer"
	default:
		panic(fmt.Errorf("unrecognized Zkind value %#v", i))
	}
}

// WriteToGo writes the zebrapack schema to w as a Go source file.
func (s *Schema) WriteToGo(w io.Writer, path string, pkg string) (err error) {
	if pkg == "" {
		fmt.Fprintf(w, "\npackage %s\n\n", s.SourcePackage)
	} else {
		fmt.Fprintf(w, "\npackage %s\n\n", pkg)
	}
	fmt.Fprintf(w, "// File re-generated by: 'zebrapack -write-to-go %s'.\n", path)
	fmt.Fprintf(w, "// The '%s' schema was originally created from: '%s'.\n\n", path, s.SourcePath)

	if len(s.Imports) > 0 {
		fmt.Fprintf(w, "import (\n")
	}
	for i := range s.Imports {
		fmt.Fprintf(w, "  %s\n", s.Imports[i])
	}
	if len(s.Imports) > 0 {
		fmt.Fprintf(w, ")\n\n")
	}

	fmt.Fprintf(w, "const zebraSchemaId64 = 0x%x // %v\n\n",
		s.ZebraSchemaId, s.ZebraSchemaId)

	for i := range s.Structs {
		err = s.Structs[i].WriteToGo(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Struct) WriteToGo(w io.Writer) (err error) {
	fmt.Fprintf(w, "\ntype %s struct {\n", s.StructName)
	for _, f := range s.Fields {
		needMsg := false
		zid := fmt.Sprintf("`zid:\"%v\"", f.Zid)
		msg := "msg:\""
		if f.FieldTagName != f.FieldGoName {
			msg += f.FieldTagName
			needMsg = true
		}
		if f.OmitEmpty {
			msg += ",omitempty"
			needMsg = true
		}
		if f.ShowZero {
			msg += ",showzero"
			needMsg = true
		}
		if f.Deprecated {
			msg += ",deprecated"
			needMsg = true
		}
		if needMsg {
			zid += " " + msg + "\""
		}
		fmt.Fprintf(w, "    %s %s %s`\n", f.FieldGoName, f.FieldTypeStr, zid)
	}
	fmt.Fprintf(w, "}\n\n")
	return nil
}

// ErrNoStructNameFound is returned by ZebraToMsgp2 when it cannot locate the
// embedded struct name string.
var ErrNoStructNameFound = fmt.Errorf("error: no -1:struct-name field:value found in zebrapack struct")

func (sch *Schema) ZebraToMsgp2(bts []byte, ignoreMissingStructName bool) (out []byte, left []byte, err error) {

	// write key:value pairs to newMap. At then end,
	// once we know how many pairs we have, then
	// we can write a map header to out and append
	// newMap after.
	//
	// We don't know the size of {the union
	// of fields present and fields absent but marked
	// showZero} until we've scanned the full bts.
	var newMap []byte

	// get the -1 key out of the map.
	var n uint32
	var nbs msgp.NilBitsStack
	n, bts, err = nbs.ReadMapHeaderBytes(bts)
	origMapFields := bts
	if err != nil {
		panic(err)
		return nil, nil, err
	}

	var fnum int
	var name string
	foundMinusOne := false

findMinusOneLoop:
	for i := uint32(0); i < n; i++ {
		fnum, bts, err = nbs.ReadIntBytes(bts)
		if fnum == -1 {
			name, bts, err = nbs.ReadStringBytes(bts)
			//fmt.Printf("\n found name = '%#v'\n", name)
			if err != nil {
				panic(err)
			}
			foundMinusOne = true
			break findMinusOneLoop
		}
		bts, err = msgp.Skip(bts)
		if err != nil {
			panic(err)
		}
	}

	if !foundMinusOne {
		if !ignoreMissingStructName {
			return nil, nil, ErrNoStructNameFound
		}
	}

	// INVAR: name is set. lookup the fields.
	tr, found := sch.Structs[name]
	if !found {
		foundMinusOne = false
	}
	// tr can be nil if we have no Schema, for example.

	// we might have more fields after adding the
	// showzero fields. write the new header after
	// we'vre the struct.
	numFieldsSeen := 0

	// translate to msgpack2
	//out = msgp.AppendMapHeader(out, n-1)

	// track found fields, do showzero
	nextFieldExpected := 0

	// re-read
	bts = origMapFields
	for i := uint32(0); i < n; i++ {
		//p("i = %v", i)
		fnum, bts, err = nbs.ReadIntBytes(bts)
		if err != nil {
			panic(err)
		}
		//p("fnum = %v", fnum)
		if fnum == -1 {
			bts, err = msgp.Skip(bts)
			if err != nil {
				panic(err)
			}
			continue
		}

		if foundMinusOne {

			// PRE: fields must arrive in sorted ascending order, in sequence,
			// monotonically increasing.
			newMap, nextFieldExpected, numFieldsSeen = zeroUpTo(tr, fnum, newMap, nextFieldExpected, numFieldsSeen)
			// encode fnum-> string translation for field name, then the field following
			newMap = msgp.AppendString(newMap, tr.Fields[fnum].FieldTagName)
			nextFieldExpected = fnum + 1
			numFieldsSeen++
		} else {
			// compensate with a fallback when no schema present:
			// just stringify the zid number so it shows up in the json.
			newMap = msgp.AppendString(newMap, fmt.Sprintf("%v", fnum))
			numFieldsSeen++
		}

		// arrays and maps need to be recursively decoded.
		newMap, bts, err = sch.zebraToMsgp2helper(bts, newMap, ignoreMissingStructName)
		if err != nil {
			panic(err)
		}
	}

	if foundMinusOne {
		// done with available fields, are any remaining in the schema?
		newMap, _, numFieldsSeen = zeroUpTo(tr, len(tr.Fields), newMap, nextFieldExpected, numFieldsSeen)
	}

	// put a header in front of the newMap pairs... now that we know
	// how many fields we have seen, so we can.
	//p("at end of ZebraToMsgp2, numFieldsSeen = %v", numFieldsSeen)
	out = msgp.AppendMapHeader(out, uint32(numFieldsSeen))
	out = append(out, newMap...)

	return out, bts, nil
}

func (sch *Schema) zebraToMsgp2helper(bts []byte, startOut []byte,
	ignoreMissingStructName bool) (out []byte, left []byte, err error) {

	out = startOut
	var nbs msgp.NilBitsStack

	k := msgp.NextType(bts)
	switch k {
	case msgp.MapType:
		// recurse
		var o2 []byte
		o2, bts, err = sch.ZebraToMsgp2(bts, ignoreMissingStructName)
		out = append(out, o2...)
	case msgp.ArrayType:
		// recurse
		var sz uint32
		sz, bts, err = nbs.ReadArrayHeaderBytes(bts)
		if err != nil {
			return nil, nil, err
		}
		out = msgp.AppendArrayHeader(out, sz)
		for i := uint32(0); i < sz; i++ {
			out, bts, err = sch.zebraToMsgp2helper(bts, out, ignoreMissingStructName)
			if err != nil {
				return nil, nil, err
			}
		}
	default:
		// find the end of the next field
		var end []byte
		end, err = msgp.Skip(bts)
		if err != nil {
			panic(err)
		}
		// copy field directly
		sz := len(bts) - len(end)
		out = append(out, bts[:sz]...)
		bts = end
	}

	return out, bts, nil
}

func writeZeroMsgpValueFor(fld *Field, out []byte) []byte {
	switch fld.FieldCategory {
	case BaseElemCat:
		switch fld.FieldPrimitive {
		case Invalid:
			panic("invalid type")
		case Bytes:
			return msgp.AppendBytes(out, []byte{})
		case String:
			return msgp.AppendString(out, "")
		case Float32, Float64, Complex64, Complex128,
			Uint, Uint8, Uint16, Uint32, Uint64,
			Byte, Int, Int8, Int16, Int32, Int64:
			return append(out, 0)
		case Bool:
			return msgp.AppendBool(out, false)
		case Intf:
			return msgp.AppendNil(out)
		case Time:
			return append(out, 0)
		case Ext:
			return msgp.AppendNil(out)
			// IDENT means an unrecognized identifier;
			// it typically means a named struct type.
			// The Str field in the Ztype will hold the
			// name of the struct.
		case IDENT:
			return msgp.AppendNil(out)
		}
	case MapCat:
		return msgp.AppendNil(out)
	case StructCat:
		return msgp.AppendNil(out)
	case SliceCat:
		return msgp.AppendNil(out)
	case ArrayCat:
		return msgp.AppendNil(out)
	case PointerCat:
		return msgp.AppendNil(out)
	}
	return msgp.AppendNil(out)
}

// zeroUpTo() starts from k and stops after stayBelow -1;
// it does nothing if k >= stayBelow. Otherwise, for
// each field, it handles the ShowZero flag: if
// the field is missing and marked ShowZero, then
// we write the field name and a zero type to
// the msgpack bytes, appending to `out`.
//
// tr cannot be nil.
func zeroUpTo(tr *Struct, stayBelow int, out []byte, k, numFieldsSeen int) (newOut []byte, newNextFieldExpected int, newFieldsSeen int) {
	//p("zeroUpTo called with k = %v, stayBelow = %v, tr = %#v", k, stayBelow, tr)
	for k < stayBelow {
		//p("k is now %v", k)
		if tr.Fields[k].Skip {
			//p("skipping field k=%v", k)
			k++
			continue
		}
		//p("tr.Fields[k] = '%#v'", tr.Fields[k])
		// fill in missing fields that are showzero
		if tr.Fields[k].ShowZero {
			numFieldsSeen++
			//p("found showzero field at k = %v", k)
			out = msgp.AppendString(out, tr.Fields[k].FieldTagName)
			out = writeZeroMsgpValueFor(&(tr.Fields[k]), out)
		}
		k++
	}
	return out, k, numFieldsSeen
}

func p(format string, args ...interface{}) {
	fmt.Printf("\n"+format+"\n", args...)
}
