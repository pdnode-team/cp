const data = document.getElementById('data')

const group1 = document.getElementById('char1-group') as HTMLFieldSetElement;
group1.disabled = true
group1.hidden = true

const characters = JSON.parse(data?.dataset?.characters ?? '[]')
console.log("characters: ", characters)

const name1 = document.getElementById('char1-name') as HTMLInputElement;

const hintUpdateText = document.getElementById('hint-update-text') as HTMLParagraphElement;

const char1Select = document.getElementById('char1-select')
char1Select?.addEventListener('change', (e) => {
    group1.hidden = false
    const value = (e.target as HTMLSelectElement).value;
    if (value == "new"){
        group1.disabled = false
        hintUpdateText.hidden = true
        return
    } else {
        hintUpdateText.hidden = false
        group1.disabled = true
    }

    name1.value = characters[Number(value) - 1].name
    
})

hintUpdateText.hidden = true

hintUpdateText?.addEventListener('click', () => {
    group1.disabled = !group1.disabled
})