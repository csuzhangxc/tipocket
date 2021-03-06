package rwregister

import (
	"regexp"
	"strconv"

	"github.com/pingcap/tipocket/pkg/elle/core"
)

var (
	opPattern = regexp.MustCompile(`([rw])([a-zA-Z])([0-9_]+)(.*)`)
)

// Int can be an int value or nil
type Int struct {
	IsNil bool
	Val   int
}

// NewInt creates Int with int value
func NewInt(v int) Int {
	return Int{
		IsNil: false,
		Val:   v,
	}
}

// NewNil creates Int with nil value
func NewNil() Int {
	return Int{
		IsNil: true,
		Val:   0,
	}
}

func (i Int) String() string {
	if i.IsNil {
		return "nil"
	}
	return strconv.Itoa(i.Val)
}

// Eq ...
func (i Int) Eq(another Int) bool {
	return i.IsNil == another.IsNil && i.Val == another.Val
}

// EqNotNil will get false for nil
func (i Int) EqNotNil(another Int) bool {
	if i.IsNil || another.IsNil {
		return false
	}
	return i.Val == another.Val
}

// IntPtr copy int and return its pointer
func IntPtr(i int) *int {
	return &i
}

// MustParseOp ...
func MustParseOp(opStr string) core.Op {
	op := core.Op{
		Type:  core.OpTypeOk,
		Value: new([]core.Mop),
	}

	for opStr != "" {
		opMatch := opPattern.FindStringSubmatch(opStr)
		if len(opMatch) != 5 {
			break
		}
		opStr = opMatch[4]
		var (
			mopType core.MopType
			mopKey  = opMatch[2]
		)
		switch opMatch[1] {
		case "r":
			mopType = core.MopTypeRead
		case "w":
			mopType = core.MopTypeWrite
		default:
			panic("unreachable")
		}
		var mopVal Int
		if opMatch[3] != "_" {
			mopValInt, err := strconv.Atoi(opMatch[3])
			if err != nil {
				panic(err)
			}
			mopVal = NewInt(mopValInt)
		} else {
			mopVal = NewNil()
		}
		*op.Value = append(*op.Value, core.Mop{
			T: mopType,
			M: map[string]interface{}{
				mopKey: mopVal,
			},
		})
	}

	return op
}

// Pair ...
func Pair(op core.Op) (core.Op, core.Op) {
	invoke := op.Copy()
	invoke.Type = core.OpTypeInvoke
	for index, mop := range *invoke.Value {
		if mop.IsRead() {
			for k := range mop.M {
				(*invoke.Value)[index].M[k] = NewNil()
			}
		}
	}
	return invoke, op
}
