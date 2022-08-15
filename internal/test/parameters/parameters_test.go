package parameters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/discord-gophers/goapi-gen/internal/testutil"
)

type testServer struct {
	array           []int32
	object          *Object
	complexObject   *ComplexObject
	passThrough     *string
	n1param         *string
	primitive       *int32
	primitiveString *string
	cookieParams    *GetCookieParams
	queryParams     *GetQueryFormParams
	headerParams    *GetHeaderParams
}

func (t *testServer) reset() Responser {
	t.array = nil
	t.object = nil
	t.complexObject = nil
	t.passThrough = nil
	t.n1param = nil
	t.primitive = nil
	t.primitiveString = nil
	t.cookieParams = nil
	t.queryParams = nil
	t.headerParams = nil
	return nil
}

// (GET /contentObject/{param})
func (t *testServer) GetContentObject(w http.ResponseWriter, r *http.Request, param ComplexObject) Responser {
	t.complexObject = &param
	return nil
}

// (GET /labelExplodeArray/{.param*})
func (t *testServer) GetLabelExplodeArray(w http.ResponseWriter, r *http.Request, param []int32) Responser {
	t.array = param
	return nil
}

// (GET /labelExplodeObject/{.param*})
func (t *testServer) GetLabelExplodeObject(w http.ResponseWriter, r *http.Request, param Object) Responser {
	t.object = &param
	return nil
}

// (GET /labelNoExplodeArray/{.param})
func (t *testServer) GetLabelNoExplodeArray(w http.ResponseWriter, r *http.Request, param []int32) Responser {
	t.array = param
	return nil
}

// (GET /labelNoExplodeObject/{.param})
func (t *testServer) GetLabelNoExplodeObject(w http.ResponseWriter, r *http.Request, param Object) Responser {
	t.object = &param
	return nil
}

// (GET /matrixExplodeArray/{.param*})
func (t *testServer) GetMatrixExplodeArray(w http.ResponseWriter, r *http.Request, param []int32) Responser {
	t.array = param
	return nil
}

// (GET /matrixExplodeObject/{.param*})
func (t *testServer) GetMatrixExplodeObject(w http.ResponseWriter, r *http.Request, param Object) Responser {
	t.object = &param
	return nil
}

// (GET /matrixNoExplodeArray/{.param})
func (t *testServer) GetMatrixNoExplodeArray(w http.ResponseWriter, r *http.Request, param []int32) Responser {
	t.array = param
	return nil
}

// (GET /matrixNoExplodeObject/{.param})
func (t *testServer) GetMatrixNoExplodeObject(w http.ResponseWriter, r *http.Request, param Object) Responser {
	t.object = &param
	return nil
}

// (GET /simpleExplodeArray/{param*})
func (t *testServer) GetSimpleExplodeArray(w http.ResponseWriter, r *http.Request, param []int32) Responser {
	t.array = param
	return nil
}

// (GET /simpleExplodeObject/{param*})
func (t *testServer) GetSimpleExplodeObject(w http.ResponseWriter, r *http.Request, param Object) Responser {
	t.object = &param
	return nil
}

// (GET /simpleNoExplodeArray/{param})
func (t *testServer) GetSimpleNoExplodeArray(w http.ResponseWriter, r *http.Request, param []int32) Responser {
	t.array = param
	return nil
}

// (GET /simpleNoExplodeObject/{param})
func (t *testServer) GetSimpleNoExplodeObject(w http.ResponseWriter, r *http.Request, param Object) Responser {
	t.object = &param
	return nil
}

// (GET /passThrough/{param})
func (t *testServer) GetPassThrough(w http.ResponseWriter, r *http.Request, param string) Responser {
	t.passThrough = &param
	return nil
}

// (GET /startingWithjNumber/{param})
func (t *testServer) GetStartingWithNumber(w http.ResponseWriter, r *http.Request, n1param string) Responser {
	t.n1param = &n1param
	return nil
}

// (GET /queryDeepObject)
func (t *testServer) GetDeepObject(w http.ResponseWriter, r *http.Request, params GetDeepObjectParams) Responser {
	t.complexObject = &params.DeepObj
	return nil
}

// (GET /simplePrimitive/{param})
func (t *testServer) GetSimplePrimitive(w http.ResponseWriter, r *http.Request, param int32) Responser {
	t.primitive = &param
	return nil
}

// (GET /queryForm)
func (t *testServer) GetQueryForm(w http.ResponseWriter, r *http.Request, params GetQueryFormParams) Responser {
	t.queryParams = &params
	if params.Ea != nil {
		t.array = params.Ea
	}
	if params.A != nil {
		t.array = params.A
	}
	if params.Eo != nil {
		t.object = params.Eo
	}
	if params.O != nil {
		t.object = params.O
	}
	if params.P != nil {
		t.primitive = params.P
	}
	if params.Ps != nil {
		t.primitiveString = params.Ps
	}
	if params.Ep != nil {
		t.primitive = params.Ep
	}
	if params.Co != nil {
		t.complexObject = params.Co
	}
	if params.N1s != nil {
		t.n1param = params.N1s
	}
	return nil
}

// (GET /header)
func (t *testServer) GetHeader(w http.ResponseWriter, r *http.Request, params GetHeaderParams) Responser {
	t.headerParams = &params
	if params.XPrimitive != nil {
		t.primitive = params.XPrimitive
	}
	if params.XPrimitiveExploded != nil {
		t.primitive = params.XPrimitiveExploded
	}
	if params.XArray != nil {
		t.array = params.XArray
	}
	if params.XArrayExploded != nil {
		t.array = params.XArrayExploded
	}
	if params.XObject != nil {
		t.object = params.XObject
	}
	if params.XObjectExploded != nil {
		t.object = params.XObjectExploded
	}
	if params.XComplexObject != nil {
		t.complexObject = params.XComplexObject
	}
	if params.N1StartingWithNumber != nil {
		t.n1param = params.N1StartingWithNumber
	}
	return nil
}

// (GET /cookie)
func (t *testServer) GetCookie(w http.ResponseWriter, r *http.Request, params GetCookieParams) Responser {
	t.cookieParams = &params
	if params.Ea != nil {
		t.array = params.Ea
	}
	if params.A != nil {
		t.array = params.A
	}
	if params.Eo != nil {
		t.object = params.Eo
	}
	if params.O != nil {
		t.object = params.O
	}
	if params.P != nil {
		t.primitive = params.P
	}
	if params.Ep != nil {
		t.primitive = params.Ep
	}
	if params.Co != nil {
		t.complexObject = params.Co
	}
	if params.N1s != nil {
		t.n1param = params.N1s
	}
	return nil
}

func TestParameterBinding(t *testing.T) {
	var ts testServer
	handler := Handler(&ts)

	expectedObject := Object{
		FirstName: "Alex",
		Role:      "admin",
	}

	expectedComplexObject := ComplexObject{
		Object:  expectedObject,
		ID:      12345,
		IsAdmin: true,
	}

	expectedArray := []int32{3, 4, 5}

	var expectedPrimitive int32 = 5

	var expectedPrimitiveString string = "123;456"

	var expectedN1Param string = "foo"

	// Check the passthrough case
	//  (GET /passThrough/{param})
	result := testutil.NewRequest().Get("/passThrough/some%20string").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	require.NotNil(t, ts.passThrough)
	assert.EqualValues(t, "some string", *ts.passThrough)
	ts.reset()

	// Check JSON marshaling of Content based parameter
	//  (GET /contentObject/{param})
	marshaledComplexObject, err := json.Marshal(expectedComplexObject)
	assert.NoError(t, err)
	q := fmt.Sprintf("/contentObject/%s", string(marshaledComplexObject))
	result = testutil.NewRequest().Get(q).GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedComplexObject, ts.complexObject)
	ts.reset()

	//  (GET /labelExplodeArray/{.param*})
	result = testutil.NewRequest().Get("/labelExplodeArray/.3.4.5").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, expectedArray, ts.array)
	ts.reset()

	//  (GET /labelExplodeObject/{.param*})
	result = testutil.NewRequest().Get("/labelExplodeObject/.role=admin.firstName=Alex").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedObject, ts.object)
	ts.reset()

	//  (GET /labelNoExplodeArray/{.param})
	result = testutil.NewRequest().Get("/labelNoExplodeArray/.3,4,5").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, expectedArray, ts.array)
	ts.reset()

	//  (GET /labelNoExplodeObject/{.param})
	result = testutil.NewRequest().Get("/labelNoExplodeObject/.role,admin,firstName,Alex").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedObject, ts.object)
	ts.reset()

	//  (GET /matrixExplodeArray/{.param*})
	uri := "/matrixExplodeArray/;id=3;id=4;id=5"
	result = testutil.NewRequest().Get(uri).GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, expectedArray, ts.array)
	ts.reset()

	//  (GET /matrixExplodeObject/{.param*})
	uri = "/matrixExplodeObject/;role=admin;firstName=Alex"
	result = testutil.NewRequest().Get(uri).GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedObject, ts.object)
	ts.reset()

	//  (GET /matrixNoExplodeArray/{.param})
	result = testutil.NewRequest().Get("/matrixNoExplodeArray/;id=3,4,5").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, expectedArray, ts.array)
	ts.reset()

	//  (GET /matrixNoExplodeObject/{.param})
	result = testutil.NewRequest().Get("/matrixNoExplodeObject/;id=role,admin,firstName,Alex").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedObject, ts.object)
	ts.reset()

	//  (GET /simpleExplodeArray/{param*})
	result = testutil.NewRequest().Get("/simpleExplodeArray/3,4,5").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, expectedArray, ts.array)
	ts.reset()

	//  (GET /simpleExplodeObject/{param*})
	result = testutil.NewRequest().Get("/simpleExplodeObject/role=admin,firstName=Alex").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedObject, ts.object)
	ts.reset()

	//  (GET /simpleNoExplodeArray/{param})
	result = testutil.NewRequest().Get("/simpleNoExplodeArray/3,4,5").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, expectedArray, ts.array)
	ts.reset()

	//  (GET /simpleNoExplodeObject/{param})
	result = testutil.NewRequest().Get("/simpleNoExplodeObject/role,admin,firstName,Alex").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedObject, ts.object)
	ts.reset()

	//  (GET /simplePrimitive/{param})
	result = testutil.NewRequest().Get("/simplePrimitive/5").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedPrimitive, ts.primitive)
	ts.reset()

	//  (GET /startingWithNumber/{1param})
	result = testutil.NewRequest().Get("/startingWithNumber/foo").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedN1Param, ts.n1param)
	ts.reset()

	// ---------------------- Test Form Query Parameters ----------------------
	//  (GET /queryForm)

	// unexploded array
	result = testutil.NewRequest().Get("/queryForm?a=3,4,5").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, expectedArray, ts.array)
	ts.reset()

	// exploded array
	result = testutil.NewRequest().Get("/queryForm?ea=3&ea=4&ea=5").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, expectedArray, ts.array)
	ts.reset()

	// unexploded object
	result = testutil.NewRequest().Get("/queryForm?o=role,admin,firstName,Alex").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedObject, ts.object)
	ts.reset()

	// exploded object
	result = testutil.NewRequest().Get("/queryForm?role=admin&firstName=Alex").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedObject, ts.object)
	ts.reset()

	// exploded primitive
	result = testutil.NewRequest().Get("/queryForm?ep=5").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedPrimitive, ts.primitive)
	ts.reset()

	// unexploded primitive
	result = testutil.NewRequest().Get("/queryForm?p=5").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedPrimitive, ts.primitive)
	ts.reset()

	// primitive string within reserved char, i.e., ';' escaped to '%3B'
	result = testutil.NewRequest().Get("/queryForm?ps=123%3B456").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedPrimitiveString, ts.primitiveString)
	ts.reset()

	// complex object
	q = fmt.Sprintf("/queryForm?co=%s", string(marshaledComplexObject))
	result = testutil.NewRequest().Get(q).GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedComplexObject, ts.complexObject)
	ts.reset()

	// starting with number
	result = testutil.NewRequest().Get("/queryForm?1s=foo").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedN1Param, ts.n1param)
	ts.reset()

	// complex object via deepObject
	do := `deepObj[Id]=12345&deepObj[IsAdmin]=true&deepObj[Object][firstName]=Alex&deepObj[Object][role]=admin`
	q = "/queryDeepObject?" + do
	result = testutil.NewRequest().Get(q).GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedComplexObject, ts.complexObject)
	ts.reset()

	// ---------------------- Test Header Query Parameters --------------------

	// unexploded header primitive.
	result = testutil.NewRequest().WithHeader("X-Primitive", "5").Get("/header").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedPrimitive, ts.primitive)
	ts.reset()

	// exploded header primitive.
	result = testutil.NewRequest().WithHeader("X-Primitive-Exploded", "5").Get("/header").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedPrimitive, ts.primitive)
	ts.reset()

	// unexploded header array
	result = testutil.NewRequest().WithHeader("X-Array", "3,4,5").Get("/header").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, expectedArray, ts.array)
	ts.reset()

	// exploded header array
	result = testutil.NewRequest().WithHeader("X-Array-Exploded", "3,4,5").Get("/header").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, expectedArray, ts.array)
	ts.reset()

	// unexploded header object
	result = testutil.NewRequest().WithHeader("X-Object",
		"role,admin,firstName,Alex").Get("/header").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedObject, ts.object)
	ts.reset()

	// exploded header object
	result = testutil.NewRequest().WithHeader("X-Object-Exploded",
		"role=admin,firstName=Alex").Get("/header").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedObject, ts.object)
	ts.reset()

	// complex object
	result = testutil.NewRequest().WithHeader("X-Complex-Object",
		string(marshaledComplexObject)).Get("/header").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedComplexObject, ts.complexObject)
	ts.reset()

	// starting with number
	result = testutil.NewRequest().WithHeader("1-Starting-With-Number",
		"foo").Get("/header").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedN1Param, ts.n1param)
	ts.reset()

	// ------------------------- Test Cookie Parameters ------------------------
	result = testutil.NewRequest().WithCookieNameValue("p", "5").Get("/cookie").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedPrimitive, ts.primitive)
	ts.reset()

	result = testutil.NewRequest().WithCookieNameValue("ep", "5").Get("/cookie").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedPrimitive, ts.primitive)
	ts.reset()

	result = testutil.NewRequest().WithCookieNameValue("a", "3,4,5").Get("/cookie").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, expectedArray, ts.array)
	ts.reset()

	result = testutil.NewRequest().WithCookieNameValue(
		"o", "role,admin,firstName,Alex").Get("/cookie").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedObject, ts.object)
	ts.reset()

	result = testutil.NewRequest().WithCookieNameValue("1s", "foo").Get("/cookie").GoWithHTTPHandler(t, handler)
	assert.Equal(t, http.StatusOK, result.Code())
	assert.EqualValues(t, &expectedN1Param, ts.n1param)
	ts.reset()
}
