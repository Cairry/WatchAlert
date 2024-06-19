package models

type SiteId string

type DetectionReq struct {
	SiteList []struct {
		ID string `json:"id"`
	} `json:"siteList"`
	Protocol string `json:"protocol"`
	Target   string `json:"target"`
}

type DetectionData struct {
	SiteId       string `json:"siteId"`
	Address      string `json:"address"`
	StatusCode   int    `json:"statusCode"`
	ResponseTime string `json:"responseTime"`
}

type DetectionSiteQuery struct {
	ID   string `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}

type DetectionSite struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Endpoint    string `json:"endpoint"`
	Description string `json:"description"`
}
