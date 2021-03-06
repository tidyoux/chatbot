package convey

import "github.com/glycerine/goconvey/convey/assertions"

var (
	ShouldEqual          = assertions.ShouldEqual
	ShouldNotEqual       = assertions.ShouldNotEqual
	ShouldAlmostEqual    = assertions.ShouldAlmostEqual
	ShouldNotAlmostEqual = assertions.ShouldNotAlmostEqual
	ShouldResemble       = assertions.ShouldResemble
	ShouldNotResemble    = assertions.ShouldNotResemble
	ShouldPointTo        = assertions.ShouldPointTo
	ShouldNotPointTo     = assertions.ShouldNotPointTo
	ShouldBeNil          = assertions.ShouldBeNil
	ShouldNotBeNil       = assertions.ShouldNotBeNil
	ShouldBeTrue         = assertions.ShouldBeTrue
	ShouldBeFalse        = assertions.ShouldBeFalse
	ShouldBeZeroValue    = assertions.ShouldBeZeroValue

	ShouldBeGreaterThan          = assertions.ShouldBeGreaterThan
	ShouldBeGreaterThanOrEqualTo = assertions.ShouldBeGreaterThanOrEqualTo
	ShouldBeLessThan             = assertions.ShouldBeLessThan
	ShouldBeLessThanOrEqualTo    = assertions.ShouldBeLessThanOrEqualTo
	ShouldBeBetween              = assertions.ShouldBeBetween
	ShouldNotBeBetween           = assertions.ShouldNotBeBetween
	ShouldBeBetweenOrEqual       = assertions.ShouldBeBetweenOrEqual
	ShouldNotBeBetweenOrEqual    = assertions.ShouldNotBeBetweenOrEqual

	ShouldContain    = assertions.ShouldContain
	ShouldNotContain = assertions.ShouldNotContain
	ShouldBeIn       = assertions.ShouldBeIn
	ShouldNotBeIn    = assertions.ShouldNotBeIn
	ShouldBeEmpty    = assertions.ShouldBeEmpty
	ShouldNotBeEmpty = assertions.ShouldNotBeEmpty

	ShouldStartWith           = assertions.ShouldStartWith
	ShouldNotStartWith        = assertions.ShouldNotStartWith
	ShouldEndWith             = assertions.ShouldEndWith
	ShouldNotEndWith          = assertions.ShouldNotEndWith
	ShouldBeBlank             = assertions.ShouldBeBlank
	ShouldNotBeBlank          = assertions.ShouldNotBeBlank
	ShouldContainSubstring    = assertions.ShouldContainSubstring
	ShouldNotContainSubstring = assertions.ShouldNotContainSubstring

	ShouldPanic        = assertions.ShouldPanic
	ShouldNotPanic     = assertions.ShouldNotPanic
	ShouldPanicWith    = assertions.ShouldPanicWith
	ShouldNotPanicWith = assertions.ShouldNotPanicWith

	ShouldHaveSameTypeAs    = assertions.ShouldHaveSameTypeAs
	ShouldNotHaveSameTypeAs = assertions.ShouldNotHaveSameTypeAs
	ShouldImplement         = assertions.ShouldImplement
	ShouldNotImplement      = assertions.ShouldNotImplement

	ShouldHappenBefore         = assertions.ShouldHappenBefore
	ShouldHappenOnOrBefore     = assertions.ShouldHappenOnOrBefore
	ShouldHappenAfter          = assertions.ShouldHappenAfter
	ShouldHappenOnOrAfter      = assertions.ShouldHappenOnOrAfter
	ShouldHappenBetween        = assertions.ShouldHappenBetween
	ShouldHappenOnOrBetween    = assertions.ShouldHappenOnOrBetween
	ShouldNotHappenOnOrBetween = assertions.ShouldNotHappenOnOrBetween
	ShouldHappenWithin         = assertions.ShouldHappenWithin
	ShouldNotHappenWithin      = assertions.ShouldNotHappenWithin
	ShouldBeChronological      = assertions.ShouldBeChronological

	ShouldBeEqualIgnoringSpaces = assertions.ShouldBeEqualIgnoringSpaces

	ShouldMatchModuloSpaces         = assertions.ShouldMatchModuloSpaces
	ShouldMatchModuloWhiteSpace     = assertions.ShouldMatchModuloWhiteSpace
	ShouldStartWithModuloWhiteSpace = assertions.ShouldStartWithModuloWhiteSpace

	// with language specific comment-until-end-of-line removal first.
	ShouldMatchModuloWhiteSpaceAndLuaComments     = assertions.ShouldMatchModuloWhiteSpaceAndLuaComments
	ShouldStartWithModuloWhiteSpaceAndLuaComments = assertions.ShouldStartWithModuloWhiteSpaceAndLuaComments

	ShouldMatchModuloWhiteSpaceAndGolangComments     = assertions.ShouldMatchModuloWhiteSpaceAndGolangComments
	ShouldStartWithModuloWhiteSpaceAndGolangComments = assertions.ShouldStartWithModuloWhiteSpaceAndGolangComments
	ShouldMatchRegex                                 = assertions.ShouldMatchRegex
)
