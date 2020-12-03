package main

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/genproto/googleapis/type/latlng"
)

// InterpretationEachValueForTime convert string to time.Time
func (f *Firestore) InterpretationEachValueForTime(mps map[string]interface{}) map[string]interface{} {
	for k, i := range mps {
		v := reflect.ValueOf(i)
		mps[k] = i // set default

		switch v.Kind() {
		case reflect.String: // for timestamp
			if tm, isOk := f.assertTimestampType(v.String()); isOk {
				mps[k] = *tm
			}
		case reflect.Map:
			vi := v.Interface()
			Debugf("type switch() map for LatLng / DocumentRef %#v", vi)
			if ms, ok := vi.(map[string]interface{}); ok {
				Debugf("map value of %#v", ms)
				// if latLng, isOk := f.assertLatLngType(ms); isOk {
				// 	Debugf("set LatLng %v, %#v", isOk, latLng)
				// 	mps[k] = latLng
				// }
				// if docRef, isOk := f.assertDocumentRef(ms); isOk {
				// 	Debugf("set DocumentRef %v, %#v", isOk, docRef)
				// 	mps[k] = docRef
				// }
				convertedMap := make(map[string]interface{})
				for mk, mv := range ms {
					if slice, isOk := f.assertSliceType(mv); isOk {
						Debugf("set key: %s value: %#v", isOk, slice)
						convertedMap[mk] = slice
					} else {
						convertedMap[mk] = mv
					}
				}
				mps[k] = convertedMap
			}
			//	case firestore.DocumentRef:
			// Debugf("type switch()  docRef %#v", v)
		case reflect.Slice:
			if slice, isOk := f.assertSliceType(i); isOk {
				Debugf("set Slice %v, %#v", isOk, slice)
				mps[k] = slice
			}
		default:
			Debugf("type switch() %#v", v)
			mps[k] = i
		}
	}
	return mps
}

func includeStringSlice(s string, ss []string) bool {
	for _, ts := range ss {
		if ts == s {
			return true
		}
	}
	return false
}

func (f *Firestore) assertDocumentRef(x map[string]interface{}) (*firestore.DocumentRef, bool) {

	hasRefKeys := []string{}
	firestorePath := ""
	if len(x) < 2 {
		return nil, false
	}
	for k, v := range x {
		if includeStringSlice(k, []string{"ID", "Path"}) {
			hasRefKeys = append(hasRefKeys, k)
			if k == "Path" {
				if path, ok := v.(string); ok {
					if doclist := strings.Split(path, "(default)/documents/"); len(doclist) >= 2 {
						firestorePath = doclist[1]
					}
				}
			}
		}
	}
	if len(hasRefKeys) != 2 || len(firestorePath) < 1 {
		Debugf("invalid hasKeys:%v, path:%v, val:%v", hasRefKeys, firestorePath, x)
		return nil, false
	}
	return f.Client.Doc(firestorePath), true
}

func (f *Firestore) assertLatLngType(x map[string]interface{}) (*latlng.LatLng, bool) {
	isOnlyLatLngKey := true
	isOnlyIntValue := true
	latLng := &latlng.LatLng{}
	if len(x) != 2 {
		return latLng, false
	}

	for k, v := range x {
		if !includeStringSlice(k, []string{"latitude", "longitude"}) {
			Debugf("invalid key %v", k)
			isOnlyLatLngKey = false
			break
		}
		if newValue, ok := v.(float64); ok {
			//	Debugf("invalid value %v", newValue)
			if k == "latitude" {
				latLng.Latitude = newValue
			} else {
				latLng.Longitude = newValue
			}
		} else {
			isOnlyIntValue = false
			break
		}
	}
	return latLng, isOnlyLatLngKey && isOnlyIntValue
}

const fsTimeLayout = "2006-01-02T15:04:05Z"

func (f *Firestore) assertTimestampType(v string) (*time.Time, bool) {
	if tm, err := time.Parse(fsTimeLayout, v); err == nil {
		return &tm, true
	}
	return nil, false
}

func (f *Firestore) assertSliceType(v interface{}) ([]interface{}, bool) {
	if reflect.ValueOf(v).Kind() != reflect.Slice {
		return nil, false
	}
	Debugf("slice before converted: %#v", v)
	var converted []interface{}
	s := reflect.ValueOf(v)
	for i := 0; i < s.Len(); i++ {
		om := f.InterpretationEachValueForTime(map[string]interface{}{
			strconv.Itoa(i): s.Index(i).Interface(),
		})
		converted = append(converted, om[strconv.Itoa(i)])
	}
	return converted, true
}
