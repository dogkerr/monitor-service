package domain

import (
	_ "dogker/lintang/monitor-service/docs"
)

// Dashboard
// @Description ini data dashboard (isinya id, owner, uid, type)
type Dashboard struct {
	// id dashboard di database
	Id string `json:"id"`
	// owner /pemilik dashboard
	Owner string `json:"owner"`
	// uid dashboard di grafana
	Uid string `json:"uid"`
	// type dashboard
	Type string `json:"type"`
}

// grafana config
type GrafanaMonitorConfig struct {
	Dashboard struct {
		Inputs []struct {
			Name        string `json:"name"`
			Label       string `json:"label"`
			Description string `json:"description"`
			Type        string `json:"type"`
			PluginID    string `json:"pluginId"`
			PluginName  string `json:"pluginName"`
		} `json:"__inputs"`
		Requires []struct {
			Type    string `json:"type"`
			ID      string `json:"id"`
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"__requires"`
		Annotations struct {
			List []struct {
				BuiltIn    int    `json:"builtIn"`
				Datasource string `json:"datasource"`
				Enable     bool   `json:"enable"`
				Hide       bool   `json:"hide"`
				IconColor  string `json:"iconColor"`
				Name       string `json:"name"`
				Type       string `json:"type"`
			} `json:"list"`
		} `json:"annotations"`
		Description  string        `json:"description"`
		Editable     bool          `json:"editable"`
		GnetID       int           `json:"gnetId"`
		GraphTooltip int           `json:"graphTooltip"`
		ID           interface{}   `json:"id"`
		Iteration    int64         `json:"iteration"`
		Links        []interface{} `json:"links"`
		Panels       []struct {
			CacheTimeout    interface{} `json:"cacheTimeout,omitempty"`
			ColorBackground bool        `json:"colorBackground,omitempty"`
			ColorValue      bool        `json:"colorValue,omitempty"`
			Colors          []string    `json:"colors,omitempty"`
			Datasource      string      `json:"datasource"`
			Decimals        int         `json:"decimals,omitempty"`
			Editable        bool        `json:"editable"`
			Error           bool        `json:"error"`
			Format          string      `json:"format,omitempty"`
			Gauge           struct {
				MaxValue         int  `json:"maxValue"`
				MinValue         int  `json:"minValue"`
				Show             bool `json:"show"`
				ThresholdLabels  bool `json:"thresholdLabels"`
				ThresholdMarkers bool `json:"thresholdMarkers"`
			} `json:"gauge,omitempty"`
			GridPos struct {
				H int `json:"h"`
				W int `json:"w"`
				X int `json:"x"`
				Y int `json:"y"`
			} `json:"gridPos"`
			Height       string        `json:"height,omitempty"`
			ID           int           `json:"id"`
			Interval     interface{}   `json:"interval,omitempty"`
			Links        []interface{} `json:"links"`
			MappingType  int           `json:"mappingType,omitempty"`
			MappingTypes []struct {
				Name  string `json:"name"`
				Value int    `json:"value"`
			} `json:"mappingTypes,omitempty"`
			MaxDataPoints   int         `json:"maxDataPoints,omitempty"`
			NullPointMode   string      `json:"nullPointMode,omitempty"`
			NullText        interface{} `json:"nullText,omitempty"`
			Postfix         string      `json:"postfix,omitempty"`
			PostfixFontSize string      `json:"postfixFontSize,omitempty"`
			Prefix          string      `json:"prefix,omitempty"`
			PrefixFontSize  string      `json:"prefixFontSize,omitempty"`
			RangeMaps       []struct {
				From string `json:"from"`
				Text string `json:"text"`
				To   string `json:"to"`
			} `json:"rangeMaps,omitempty"`
			Sparkline struct {
				FillColor string `json:"fillColor"`
				Full      bool   `json:"full"`
				LineColor string `json:"lineColor"`
				Show      bool   `json:"show"`
			} `json:"sparkline,omitempty"`
			TableColumn string `json:"tableColumn,omitempty"`
			Targets     []struct {
				Expr           string `json:"expr"`
				Format         string `json:"format"`
				Hide           bool   `json:"hide"`
				IntervalFactor int    `json:"intervalFactor"`
				LegendFormat   string `json:"legendFormat"`
				RefID          string `json:"refId"`
				Step           int    `json:"step"`
			} `json:"targets"`
			// Thresholds    string `json:"thresholds,omitempty"` // bikin error ajg
			Title         string `json:"title"`
			Type          string `json:"type"`
			ValueFontSize string `json:"valueFontSize,omitempty"`
			ValueMaps     []struct {
				Op    string `json:"op"`
				Text  string `json:"text"`
				Value string `json:"value"`
			} `json:"valueMaps,omitempty"`
			ValueName   string `json:"valueName,omitempty"`
			AliasColors struct {
				SENT string `json:"SENT"`
			} `json:"aliasColors,omitempty"`
			Bars       bool `json:"bars,omitempty"`
			DashLength int  `json:"dashLength,omitempty"`
			Dashes     bool `json:"dashes,omitempty"`
			Fill       int  `json:"fill,omitempty"`
			Grid       struct {
			} `json:"grid,omitempty"`
			Legend struct {
				Avg     bool `json:"avg"`
				Current bool `json:"current"`
				Max     bool `json:"max"`
				Min     bool `json:"min"`
				Show    bool `json:"show"`
				Total   bool `json:"total"`
				Values  bool `json:"values"`
			} `json:"legend,omitempty"`
			Lines           bool          `json:"lines,omitempty"`
			Linewidth       int           `json:"linewidth,omitempty"`
			Percentage      bool          `json:"percentage,omitempty"`
			Pointradius     int           `json:"pointradius,omitempty"`
			Points          bool          `json:"points,omitempty"`
			Renderer        string        `json:"renderer,omitempty"`
			SeriesOverrides []interface{} `json:"seriesOverrides,omitempty"`
			SpaceLength     int           `json:"spaceLength,omitempty"`
			Stack           bool          `json:"stack,omitempty"`
			SteppedLine     bool          `json:"steppedLine,omitempty"`
			TimeFrom        interface{}   `json:"timeFrom,omitempty"`
			TimeShift       interface{}   `json:"timeShift,omitempty"`
			Tooltip         struct {
				MsResolution bool   `json:"msResolution"`
				Shared       bool   `json:"shared"`
				Sort         int    `json:"sort"`
				ValueType    string `json:"value_type"`
			} `json:"tooltip,omitempty"`
			Transparent bool `json:"transparent,omitempty"`
			Xaxis       struct {
				Buckets interface{}   `json:"buckets"`
				Mode    string        `json:"mode"`
				Name    interface{}   `json:"name"`
				Show    bool          `json:"show"`
				Values  []interface{} `json:"values"`
			} `json:"xaxis,omitempty"`
			Yaxes []struct {
				Format  string      `json:"format"`
				Label   interface{} `json:"label"`
				LogBase int         `json:"logBase"`
				Max     interface{} `json:"max"`
				Min     interface{} `json:"min"`
				Show    bool        `json:"show"`
			} `json:"yaxes,omitempty"`
			Yaxis struct {
				Align      bool        `json:"align"`
				AlignLevel interface{} `json:"alignLevel"`
			} `json:"yaxis,omitempty"`
			Alert struct {
				Conditions []struct {
					Evaluator struct {
						Params []float64 `json:"params"`
						Type   string    `json:"type"`
					} `json:"evaluator"`
					Query struct {
						Params []string `json:"params"`
					} `json:"query"`
					Reducer struct {
						Params []interface{} `json:"params"`
						Type   string        `json:"type"`
					} `json:"reducer"`
					Type string `json:"type"`
				} `json:"conditions"`
				ExecutionErrorState string `json:"executionErrorState"`
				Frequency           string `json:"frequency"`
				Handler             int    `json:"handler"`
				Name                string `json:"name"`
				NoDataState         string `json:"noDataState"`
				Notifications       []struct {
					ID int `json:"id"`
				} `json:"notifications"`
			} `json:"alert,omitempty"`
			Columns []struct {
				Text  string `json:"text"`
				Value string `json:"value"`
			} `json:"columns,omitempty"`
			FontSize   string      `json:"fontSize,omitempty"`
			PageSize   interface{} `json:"pageSize,omitempty"`
			Scroll     bool        `json:"scroll,omitempty"`
			ShowHeader bool        `json:"showHeader,omitempty"`
			Sort       struct {
				Col  int  `json:"col"`
				Desc bool `json:"desc"`
			} `json:"sort,omitempty"`
			Styles []struct {
				ColorMode  interface{} `json:"colorMode"`
				Colors     []string    `json:"colors"`
				Decimals   int         `json:"decimals"`
				Pattern    string      `json:"pattern"`
				Thresholds []string    `json:"thresholds"`
				Type       string      `json:"type"`
				Unit       string      `json:"unit"`
			} `json:"styles,omitempty"`
			Transform string `json:"transform,omitempty"`
		} `json:"panels"`
		Refresh       string        `json:"refresh"`
		SchemaVersion int           `json:"schemaVersion"`
		Style         string        `json:"style"`
		Tags          []interface{} `json:"tags"`
		Templating    struct {
			List []struct {
				AllValue string `json:"allValue,omitempty"`
				Current  struct {
				} `json:"current"`
				Datasource     string        `json:"datasource"`
				Hide           int           `json:"hide"`
				IncludeAll     bool          `json:"includeAll"`
				Label          string        `json:"label"`
				Multi          bool          `json:"multi"`
				Name           string        `json:"name"`
				Options        []interface{} `json:"options"`
				Query          string        `json:"query"`
				Refresh        int           `json:"refresh"`
				Regex          string        `json:"regex,omitempty"`
				Sort           int           `json:"sort,omitempty"`
				TagValuesQuery interface{}   `json:"tagValuesQuery,omitempty"`
				Tags           []interface{} `json:"tags,omitempty"`
				TagsQuery      interface{}   `json:"tagsQuery,omitempty"`
				Type           string        `json:"type"`
				UseTags        bool          `json:"useTags,omitempty"`
				Auto           bool          `json:"auto,omitempty"`
				AutoCount      int           `json:"auto_count,omitempty"`
				AutoMin        string        `json:"auto_min,omitempty"`
			} `json:"list"`
		} `json:"templating"`
		Time struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"time"`
		Timepicker struct {
			RefreshIntervals []string `json:"refresh_intervals"`
			TimeOptions      []string `json:"time_options"`
		} `json:"timepicker"`
		Timezone string `json:"timezone"`
		Title    string `json:"title"`
		UID      string `json:"uid"`
		Version  int    `json:"version"`
	} `json:"dashboard"`
}

type GrafanaLogsDashboard struct {
	Dashboard struct {
		Annotations struct {
			List []struct {
				BuiltIn    int `json:"builtIn"`
				Datasource struct {
					Type string `json:"type"`
					UID  string `json:"uid"`
				} `json:"datasource"`
				Enable    bool   `json:"enable"`
				Hide      bool   `json:"hide"`
				IconColor string `json:"iconColor"`
				Name      string `json:"name"`
				Type      string `json:"type"`
			} `json:"list"`
		} `json:"annotations"`
		Editable             bool          `json:"editable"`
		FiscalYearStartMonth int           `json:"fiscalYearStartMonth"`
		GraphTooltip         int           `json:"graphTooltip"`
		ID                   int           `json:"id"`
		Links                []interface{} `json:"links"`
		LiveNow              bool          `json:"liveNow"`
		Panels               []struct {
			Datasource struct {
				Type string `json:"type"`
				UID  string `json:"uid"`
			} `json:"datasource"`
			GridPos struct {
				H int `json:"h"`
				W int `json:"w"`
				X int `json:"x"`
				Y int `json:"y"`
			} `json:"gridPos"`
			ID      int `json:"id"`
			Options struct {
				DedupStrategy      string `json:"dedupStrategy"`
				EnableLogDetails   bool   `json:"enableLogDetails"`
				PrettifyLogMessage bool   `json:"prettifyLogMessage"`
				ShowCommonLabels   bool   `json:"showCommonLabels"`
				ShowLabels         bool   `json:"showLabels"`
				ShowTime           bool   `json:"showTime"`
				SortOrder          string `json:"sortOrder"`
				WrapLogMessage     bool   `json:"wrapLogMessage"`
			} `json:"options"`
			Targets []struct {
				Datasource struct {
					Type string `json:"type"`
					UID  string `json:"uid"`
				} `json:"datasource"`
				EditorMode string `json:"editorMode"`
				Expr       string `json:"expr"`
				QueryType  string `json:"queryType"`
				RefID      string `json:"refId"`
			} `json:"targets"`
			Title         string `json:"title"`
			Type          string `json:"type"`
			Description   string `json:"description,omitempty"`
			PluginVersion string `json:"pluginVersion,omitempty"`
		} `json:"panels"`
		Refresh       string        `json:"refresh"`
		SchemaVersion int           `json:"schemaVersion"`
		Tags          []interface{} `json:"tags"`
		Templating    struct {
			List []struct {
				Current struct {
					Selected bool   `json:"selected"`
					Text     string `json:"text"`
					Value    string `json:"value"`
				} `json:"current"`
				Hide    int    `json:"hide"`
				Name    string `json:"name"`
				Options []struct {
					Selected bool   `json:"selected"`
					Text     string `json:"text"`
					Value    string `json:"value"`
				} `json:"options"`
				Query       interface{} `json:"query"`
				SkipURLSync bool        `json:"skipUrlSync"`
				Type        string      `json:"type"`
				Description string      `json:"description,omitempty"`
				IncludeAll  bool        `json:"includeAll,omitempty"`
				Label       string      `json:"label,omitempty"`
				Multi       bool        `json:"multi,omitempty"`
				QueryValue  string      `json:"queryValue,omitempty"`
				Datasource  struct {
					Type string `json:"type"`
					UID  string `json:"uid"`
				} `json:"datasource,omitempty"`
				Definition string `json:"definition,omitempty"`
				Refresh    int    `json:"refresh,omitempty"`
				Regex      string `json:"regex,omitempty"`
				Sort       int    `json:"sort,omitempty"`
			} `json:"list"`
		} `json:"templating"`
		Time struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"time"`
		Timepicker struct {
		} `json:"timepicker"`
		Timezone  string `json:"timezone"`
		Title     string `json:"title"`
		UID       string `json:"uid"`
		Version   int    `json:"version"`
		WeekStart string `json:"weekStart"`
		Style     string `json:"style"`
	} `json:"dashboard"`
}
