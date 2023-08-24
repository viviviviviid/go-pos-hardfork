package person

import "fmt"

type Person struct { // 외부로 export하고 싶으면 대문자로 시작
	name string // 이것도 대문자가아니라서 export가 안됨
	age  int
}

func (p *Person) SetDetails(name string, age int) {
	// 여기서 *이 없으면, 복사해서 가져온 복사본을 수정하는 것이고
	// 있으면 연결된 진짜 내용을 수정
	p.name = name
	p.age = age
	fmt.Println("see details", p)
}

// 크기가 작거나, 수정할 필요가 없으면 * 을 안써도 됨
// 그 외적으로는 receiver pointer function 형태를 만들어주기 위해 * 을 사용하기
