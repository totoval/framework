package model


type Modeller interface {
	Default() interface{}
	ObjArr(filterArr [][]interface{}, sortArr []Sort, limit int, withTrashed bool) []interface{} //@todo     public function getObjArr(?array $filter_arr = [], ?array $sort_arr = null, ?int $limit = null, bool $with_trashed = false): Collection;
	ObjArrPaginate(filterArr [][]interface{}, sortArr []Sort, limit int, withTrashed bool)       //@todo     public function getObjArrPaginate(int $per_page, ?array $filter_arr = [], ?array $sort_arr = null, bool $with_trashed = false): LengthAwarePaginator;
}


//func (m *Model) shouldInstantiate() { //     private function shouldInstantiate(bool $should, $primary_key_variable = null)
//
//}
//
//func (m *Model) readOnlyGuardian() { //     private function readOnlyGuardian()
//
//}
