package starwars

type QueryImpl struct {
}

func (this *QueryImpl) Hero(episode Episode) (ch *Character) {
	return
}
func (this *QueryImpl) Human(id string) (hu *Human) {
	return
}
func (this *QueryImpl) Droid(id string) (dr *Droid) {
	return
}

var DefaultQuery Query = &QueryImpl{}
