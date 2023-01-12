var xhr = new XMLHttpRequest()
var resp = ''
var content = document.getElementById('')
var type = document.getElementById('')
var option = document.getElementById('')
var input = document.getElementById('')
var chapterList = []
var chapterIndex = 0
window.onload = function() {
    content = document.getElementById('content')
    type = document.getElementById('type')
    option = document.getElementById('option')
    input = document.getElementById('keyInput')
    let btn = document.getElementsByTagName('svg')[0]
    input.onmouseover = function() {
        input.style.filter = 'brightness(100%)'
    }

    input.onmouseout = function() {
        input.style.filter = 'brightness(80%)'
    }
    btn.onclick = searchFunc
    btn.onmouseover = function() {
        btn.style.color = 'white'
    }

    btn.onmouseout = function() {
        btn.style.color = ''
    }
    let selects = document.getElementsByTagName('select')
    for (let i in selects) {
        selects[i].onmouseover = function() {
            selects[i].style.filter = 'brightness(80%)'
        }
        selects[i].onmouseout = function() {
            selects[i].style.filter = 'brightness(100%)'
        }
    }
}

function query(url, callBack) {
    xhr.open('get', url, true)
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    xhr.send();
    xhr.onreadystatechange = function() {
        if (xhr.readyState == 4 && xhr.status == 200) {
            resp = xhr.responseText
            while (content.childNodes.length > 0) {
                content.removeChild(content.childNodes.item(0))
            }
            console.log(url)
            console.log(JSON.parse(resp))
            callBack()
        }
    }
}

function searchFunc() {
    switch(type.options[type.selectedIndex].value){
        case 'fiction':
        query('https://api.pingcc.cn/fiction/search/'+option.options[option.selectedIndex].value+'/'+input.value+'/1/30', queryFictions)
        break
        case 'comic':
        query('https://api.pingcc.cn/comic/search/'+option.options[option.selectedIndex].value+'/'+input.value+'/1/30', queryComics)
        break
        case 'video':
        query('https://api.pingcc.cn/video/search/'+option.options[option.selectedIndex].value+'/'+input.value+'/1/30', queryVideos)
        break
    }
}

function queryType(type) {
    let body = JSON.parse(resp)
    let list = document.createElement('ul')
    content.appendChild(list)
    for (let i in body.data) {
        let data = body.data[i]
        let img = document.createElement('img')
        img.src = data.cover
        let title = document.createElement('p')
        title.innerHTML = data.title
        let text = document.createElement('h5')
        let li = document.createElement('li')
        li.appendChild(img)
        li.appendChild(title)
        li.appendChild(text)
        list.appendChild(li)
        text.innerHTML = '作者: ' +data.author + '<br/>分类: ' + data.fictionType + '<br/>简介: ' + data.descs
        text.hidden = true
        text.style.position = 'absolute'
        text.style.color = 'white'
        text.style.top = img.offsetTop+"px"
        text.style.width = img.offsetWidth+"px"
        li.onmouseover = function() {
            img.style.filter = 'brightness(50%)'
            text.hidden = false
        }
        li.onmouseout = function() {
            img.style.filter = 'brightness(100%)'
            text.hidden = true
        }
        switch (type){
            case 'fiction':{
                li.onclick = function() {
                    query('https://api.pingcc.cn/fictionChapter/search/'+data.fictionId, queryFictionChapter)
                }
                break
            }
            case 'comic':{
                li.onclick = function() {
                    query('https://api.pingcc.cn/comicChapter/search/'+data.comicId, queryComicChapter)
                }
                break
            }
        }
    }
}

function queryTypeChapter(type) {
    chapterList = JSON.parse(resp).data.chapterList
    for (let i in chapterList) {
        let data = chapterList[i]
        let title = document.createElement('button')
        title.innerText = data.title
        title.onmouseover = function() {
            title.style.filter = 'brightness(80%)'
        }
        title.onmouseout = function() {
            title.style.filter = 'brightness(100%)'
        }
        switch (type) {
            case 'fiction':
                title.onclick = function() {
                    chapterIndex = i
                    query('https://api.pingcc.cn/fictionContent/search/'+data.chapterId, queryFictionData)
                }
                break
            case 'comic':
                title.onclick = function() {
                    chapterIndex = i
                    query('https://api.pingcc.cn/comicContent/search/'+data.chapterId, queryComicData)
                }
                break
        }
        content.appendChild(title)
    }
}

function queryTypeData(type) {
    content.innerHTML = '<h1>'+chapterList[chapterIndex].title+'</h1><br>'
    let data = JSON.parse(resp).data
    switch (type) {
        case 'fiction':
            for (let i in data){
                let p = document.createElement('p')
                p.innerHTML = data[i]
                content.appendChild(p)
            }
            break
        case 'comic':
            for (let i in data) {
                let img = document.createElement('img')
                img.src = data[i]
                content.appendChild(img)
            }
            break
    }
    if (chapterList.length > 1) {
        let span = document.createElement('button')
        span.style = "class='switchChapter'"
        if (chapterList.length - 1 == chapterIndex) {
            span.style.margin = '0 auto'
            span.innerHTML = '上一章'
            span.onclick = function() {
                lastChapter(type)
            }
        } else if (chapterIndex == 0) {
            span.style.margin = '0 auto'
            span.innerHTML = '下一章'
            span.onclick = function() {
                nextChapter(type)
            }
        } else {
            let span1 = document.createElement('button')
            span1.innerHTML = '上一章'
            span1.onclick = function() {
                lastChapter(type)
            }
            content.appendChild(span1)
            span.style.float = 'right'
            span.innerHTML = '下一章'
            span.onclick = function() {
                nextChapter(type)
            }
        }
        content.appendChild(span)
    }
}

function lastChapter(type) {
    chapterIndex--
    switch (type) {
        case 'fiction':
            query('https://api.pingcc.cn/fictionContent/search/'+chapterList[chapterIndex].chapterId, queryFictionData)
            break
        case 'comic':
            query('https://api.pingcc.cn/comicContent/search/'+chapterList[chapterIndex].chapterId, queryComicData)
            break
    }
}

function nextChapter(type) {
    chapterIndex--
    switch (type) {
        case 'fiction':
            query('https://api.pingcc.cn/fictionContent/search/'+chapterList[chapterIndex].chapterId, queryFictionData)
            break
        case 'comic':
            query('https://api.pingcc.cn/comicContent/search/'+chapterList[chapterIndex].chapterId, queryComicData)
            break
    }
}

function queryFictions() {
    queryType('fiction')
}

function queryComics() {
    queryType('comic')
}

function queryFictionChapter() {
    queryTypeChapter('fiction')
}

function queryComicChapter() {
    queryTypeChapter('comic')
}

function queryFictionData() {
    queryTypeData('fiction')
}

function queryComicData() {
    queryTypeData('comic')
}