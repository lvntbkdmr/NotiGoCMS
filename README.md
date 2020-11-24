![CRON Update Job](https://github.com/lvntbkdmr/NotiGoCMS/workflows/CRON%20Update%20Job/badge.svg?branch=master)
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![LinkedIn][linkedin-shield]][linkedin-url]

<!-- PROJECT LOGO -->
<br />
<p align="center">

  <h3 align="center">NotiGoCMS</h3>

  <p align="center">
    A Notion.so Based Hugo Blog Content Management System !
    <br />
    <a href="https://www.lvnt.be/">View Demo</a>
    ·
    <a href="https://github.com/lvntbkdmr/NotiGoCMS/issues">Report Bug</a>
    ·
    <a href="https://github.com/lvntbkdmr/NotiGoCMS/issues">Request Feature</a>
    ·
    <a href="https://github.com/lvntbkdmr/NotiGoCMS/pulls">Send a Pull Request</a>
  </p>
</p>



<!-- TABLE OF CONTENTS -->
## Table of Contents

* [About the Project](#about-the-project)
  * [Built With](#built-with)
* [Getting Started](#getting-started)
  * [Prerequisites](#prerequisites)
  * [Installation](#installation)
* [Usage](#usage)
* [Roadmap](#roadmap)
* [Contributing](#contributing)
* [License](#license)
* [Contact](#contact)
* [Acknowledgements](#acknowledgements)



<!-- ABOUT THE PROJECT -->
## About The Project

[![Product Name Screen Shot][product-screenshot]](https://example.com)

This program is responsible for fetching your Notion pages and converting them into markdown files along with their static contents
Converted markdown files and their static contents are then pushed to your original Hugo blog git repo
Then, you only need to trigger "hugo" command to rebuild the required html files (If netlify is being used, this is done automatically)

[kjk/notionapi](https://github.com/kjk/notionapi) is being used as api to fetch Notion data and [kjk/blog](https://github.com/kjk/blog) is being referenced for conversion

### Built With

* [Go](https://golang.org/)
* [NotionApi](https://github.com/kjk/notionapi)

<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running follow these simple example steps.

### Prerequisites

* Go
* Hugo based blog in GitHub repo
* A Notion.so page which lists the articles for your blog

### Installation & Usage

1. Clone the repo
```sh
git clone https://github.com/lvntbkdmr/NotiGoCMS.git
```
2. Edit config.toml file accordingly
```toml
[notion]
startPage = "1087c264ef6b450ca5c7d3c034b399d7"
#It will be get from your notion page link. For example my articles list page => https://www.notion.so/lvntbkdmr/Articles-1087c264ef6b450ca5c7d3c034b399d7 so my startPage should be equal to 1087c264ef6b450ca5c7d3c034b399d7.

[cms]
cacheDir = "cache" # This is the cache folder to keep your content cached inside the repo
postsDir = "blog/content/posts" # Your Hugo blog will be added as a sub-module in this repo, this is the relative path for your posts directory in your blog
imgDir = "blog/static/img" # This is the relative path to store your static images fetched from Notion servers
```
3. Edit setup_env.sh
```sh
1. Edit git configurations
2. Replace https://github.com/lvntbkdmr/blog with your Hugo blog repo link
```
4. Set your environment's **GH_TOKEN_OVERRIDE** => Generated from GitHub -> Your Profile -> Developer Settings -> Personal Access Tokens
```
Note that on GitHub Action environment, you can also use default GITHUB_ACTION env variable, 
but it does not have write access permission for forked repos as this project forkes blog 
module and pushes back the changes. In that case, you have to define a new env. variable 
in GitHub other than GITHUB_ACTION (thats why I named as GH_TOKEN_OVERRIDE) in order to 
use GitHub action jobs for "automatic" updates for your blog as I did in here.
Ref. = https://github.com/ad-m/github-push-action/issues/40
```
5. Set your environment's **NOTION_TOKEN** => 
```
To access non-public pages you need to find out authentication token.
Auth token is the value of token_v2 cookie.
In Chrome: open developer tools (Menu More Tools\Developer Tools), navigate to Application tab, 
look under Storage \ Cookies and copy the value of token_v2 cookie. You can do similar things 
in other browsers.
```
6. Build the project
```sh
go build -o NotiGoCMS
```
7. Run shell script
```sh
sh setup_env.sh
```

## Additional
You can also automate the fetching and updating process through CI tools such as GitHub Actions. I personally follow that approach, my crob job description can be found on [.github/workflows/](https://github.com/lvntbkdmr/NotiGoCMS/tree/master/.github/workflows)

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **extremely appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request



<!-- LICENSE -->
## License

TODO


<!-- CONTACT -->
## Contact

Levent BEKDEMİR - lvntbkdmr@gmail.com

Project Link: [https://github.com/lvntbkdmr/NotiGoCMS](https://github.com/lvntbkdmr/NotiGoCMS)



<!-- ACKNOWLEDGEMENTS -->
## Acknowledgements
* [Krzysztof Kowalczyk](https://github.com/kjk)



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[forks-shield]: https://img.shields.io/github/forks/lvntbkdmr/NotiGoCMS?style=for-the-badge
[forks-url]: https://github.com/lvntbkdmr/NotiGoCMS/network/members
[stars-shield]: https://img.shields.io/github/stars/lvntbkdmr/NotiGoCMS?style=for-the-badge
[stars-url]: https://github.com/lvntbkdmr/NotiGoCMS/stargazers
[issues-shield]: https://img.shields.io/github/issues/lvntbkdmr/NotiGoCMS?style=for-the-badge
[issues-url]: https://github.com/lvntbkdmr/NotiGoCMS/issues
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=flat-square&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/leventbekdemir
