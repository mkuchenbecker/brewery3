package interview

import (

	"testing"
  "github.com/stretchr/testify/assert"
)

func  TestInterviewSolution( t *testing.T)  {

	t.Run("base case 1", func(t *testing.T){
		tuples := []IPTuple{
			{
				start: 5,
				end: 30,
			},
			{
				start: 10,
				end: 25,
			},
		}
		chargeTime := GetChargeTime(tuples,0,40)
		assert.Equal(t,int64(15),chargeTime)
	})

}
