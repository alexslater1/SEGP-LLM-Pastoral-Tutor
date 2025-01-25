package tools

import (
	"os"
	"testing"

	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/stretchr/testify/assert"
)

func TestParseGoogleSearchResults(t *testing.T) {
	html := `<div class="MjjYud"><div class="A6K0A" data-rpos="20" id="wxVGi"><div jscontroller="SC7lYd" class="g Ww4FFb vt6azd tF2Cxc asEBEc" lang="en" style="width:inherit" jsaction="QyLbLe:OMITjf;ewaord:qsYrDe;xd28Mb:A6j43c" data-hveid="CGEQAA" data-ved="2ahUKEwip_pmJofqKAxW4WEEAHSHHMHEQFSgAegQIYRAA"><div class="srKDX" data-snc="BudzWd"><div class="kb0PBd ieodic jGGQ5e" style="grid-area:x5WNvb" data-snf="x5WNvb" data-snhf="0"><div class="yuRUbf"><div><span jscontroller="msmzHf" jsaction="rcuQ6b:npT2md;PYDNKe:bLV6Bd;mLt3mc"><a jsname="UWckNb" href="https://en.wikipedia.org/wiki/Hello" data-ved="2ahUKEwip_pmJofqKAxW4WEEAHSHHMHEQFnoECE8QAQ" ping="/url?sa=t&amp;source=web&amp;rct=j&amp;opi=89978449&amp;url=https://en.wikipedia.org/wiki/Hello&amp;ved=2ahUKEwip_pmJofqKAxW4WEEAHSHHMHEQFnoECE8QAQ"><br><h3 class="LC20lb MBeuO DKV0Md">Hello</h3><div class="notranslate HGLrXd NJjxre iUh30 ojE3Fb"><div class="q0vns"><span class="DDKf1c"><div class="eqA2re UnOTSe Vwoesf" aria-hidden="true"><img class="XNo5Ab" src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABwAAAAcCAAAAABXZoBIAAAAnklEQVR4AeTNIQiDQABG4b+u17X1aF6PK3YEO9iMJqPVau82y4FgMezS0oVLhqsHtrcqeqzDXv3CEz/6L4yTtZM3dnHmPTtjzXZAXKYVo4agkU2GI2Lloc6JDez1+flswMu1EQZ3xlE7lK8eKDkjtwE+crBMV+wesKmCiisGGepZIfQJpMj9SNb2MYWrChjVkULuCyCfRvsdmBieyQQAsoDk/9ryhFMAAAAASUVORK5CYII=" style="height:26px;width:26px" alt="" data-csiid="s_2IZ-nYB7ixhbIPoY7DiQc_13" data-atf="4"></div></span><div class="CA5RN"><div><span class="VuuXrf">Wikipedia</span></div><div class="byrV5b"><cite class="qLRx3b tjvcx GvPZzd cHaqb" role="text">https://en.wikipedia.org<span class="ylgVCe ob9lvb" role="text"> › wiki › Hello</span></cite></div></div></div></div><span jscontroller="IX53Tb" jsaction="rcuQ6b:npT2md" style="display:none"></span></a></span><div class="B6fmyf byrV5b Mg1HEd"><div class="HGLrXd iUh30 ojE3Fb"><div class="q0vns"><span class="DDKf1c"><div class="eqA2re UnOTSe" style="height:26px;width:26px"></div></span><div class="CA5RN"><div><span class="VuuXrf">Wikipedia</span></div><div class="byrV5b"><cite class="qLRx3b tjvcx GvPZzd cHaqb" role="text">https://en.wikipedia.org<span class="ylgVCe ob9lvb" role="text"> › wiki › Hello</span></cite></div></div></div></div><div class="csDOgf BCF2pd ezY6nb L48a4c"><div jscontroller="gOTY1" data-id="atritem-https://en.wikipedia.org/wiki/Hello" jsdata="PFrTzf;_;BtQE/Q" data-viewer-group="1" jsaction="rcuQ6b:npT2md;aevozb:T2P31d;vcOT6c:C6KsF;k7WJpc:beCLof;jH1Skf:sCDZjb"><div><div jsdata="l7Bhpb;_;BtQFBM" jscontroller="PbHo4e" jsshadow="" jsaction="rcuQ6b:npT2md;h5M12e;jGQF0b:kNqZ1c;" data-viewer-entrypoint="1" data-ved="2ahUKEwip_pmJofqKAxW4WEEAHSHHMHEQ2esEegQITxAJ"><div jsslot=""><div jsname="I3kE2c" class="MJ8UF iTPLzd rNSxBe eY4mx lUn2nc" style="position:absolute" aria-label="About this result" role="button" tabindex="0"><span jsname="czHhOd" class="D6lY4c mBswFe"><span jsname="Bil8Ae" class="xTFaxe z1asCe" style="height:18px;line-height:18px;width:18px"><svg focusable="false" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12 8c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2zm0 2c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm0 6c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2z"></path></svg></span></span></div></div></div></div></div></div></div></div></div></div><div class="kb0PBd LnCrMe" style="grid-area:Vjbam;width:92px;padding-left:20px" data-snf="Vjbam" data-sncf="0,1,2,3"><div><a href="https://en.wikipedia.org/wiki/Hello" data-ved="2ahUKEwip_pmJofqKAxW4WEEAHSHHMHEQqa4BegQIVxAA" ping="/url?sa=t&amp;source=web&amp;rct=j&amp;opi=89978449&amp;url=https://en.wikipedia.org/wiki/Hello&amp;ved=2ahUKEwip_pmJofqKAxW4WEEAHSHHMHEQqa4BegQIVxAA"><div style="border-radius:8px;width:92px;height:92px" class="uhHOwf ez24Df"><img alt="hello from en.wikipedia.org" id="dimg_s_2IZ-nYB7ixhbIPoY7DiQc_19" src="data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wCEAAkGBwgHBgkIBwgKCgkLDRYPDQwMDRsUFRAWIB0iIiAdHx8kKDQsJCYxJx8fLT0tMTU3Ojo6Iys/RD84QzQ5OjcBCgoKDQwNGg8PGjclHyU3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3Nzc3N//AABEIAFwAXAMBIgACEQEDEQH/xAAcAAADAQEBAQEBAAAAAAAAAAAFBgcEAwIAAQj/xABEEAABAwIEAwQFBwcNAAAAAAABAgMRAAQFEiExE0FRBiJhcRQygZGxByMzQlKh8BYkkqKzwdElNDZDU2JjcnOCg6Ph/8QAGAEBAQEBAQAAAAAAAAAAAAAAAAECBAP/xAAZEQEBAAMBAAAAAAAAAAAAAAAAAQIRMSH/2gAMAwEAAhEDEQA/ALJdKWm3dU2QFhCiknaYoU3eXqUjiOtlXQIgeyt2IDiWV0gAkqZWB7jQdHdSAdDG451yxttF/czPGQEx9iuwurjLJeT+hQ8kFWU7kQdd66EgJEGPKqNvpNwY76Rr9nWvhcXEauJ/QrEhZBGp360oYv2suziT9vhvCbaaJSHFDMXCN+UATpQPSn7kE/PtwD9ivSLi4Mkuog7d2Km+Hdubw4ozZYoyzw3nEoS+13Skk90lMxG3TrrFPYUUmCYPSmhv49wDAca9qf8A2vi/cQYW0T1yn+NYEO96JPKuvEgzQErYuluXikqk6pEeVdK42KipjMd5PxrvNZqs99IsrjLoeEqPcalB+UJ7gtrNo3Ck6Qg+Ebq1FVi9/mVxrrwlfA1/PzbSF2ttrEwZ5gZRMVucQ/4ZjruL4LeXBBbLL1ulOWQRLgnY9KUbXtPitoHou33EpdITxXc8CTyVIH3Uz9l2m28AxJOyRcWxMDlmG330g4gUo4+Vw/Snl570DjgXbDF8Zu2LC3Zt/SbgkIccRATAJJMGIgHl4V7T2XxO9uHeHbWyFoUpDq0PONgHXrmBPs50N+S+xuHsdauUNK9GaQ4HlFHcyqSoBM8zJGnSaq902VsOJafLCj6rqQJB2G+lBH8a7P4nhd4h26BSht1C0uzKdFA7idNtwK3v45jSmH7hWJKGVMgN9zw3EfgU39qEqQhrDkFVyq4RBKl9/MFCZ1nmnQddopPxm1ebt7hl5stuED6XSBI18vfQb7btNinCbKrholSEn1BOu/wr0cdxZevpqAnw05edCbS3R6OzJI7iYkGfCda7qShMgqzCRqE1UUrsW4892et13DpedK3ApZMzC1R90UcigPYQJHZm3ymQVun9dVHtayrNen8xuP8ASX8DUDQgJs2IVr5bfgir5f62FyNpZX8DUEQ0U2rBJSITrprGnQ1ZwNPZ1Ya7N4vlIJTdWs6/4iaWUdnb27zquCbS24hWCtErUJ5InygmB50dwVxLeBYqhSQeI6ypKY5gkjfxTMbaVmvXVgIWkLU4DnKgqMwjX3ztVBvAcYGGNW9vat5GG+6lon1hrOYx625zdSeWlMWNdprKysGblvK+HzlS2TljUAzIgROxqbXDgu7EXDQgyCYGVQ1mZHSv1+6dRaPlo9/hgupjcAbDyOvT1utAwovS5eLK3nLXN3UusLhTfMQd4nfrJ8QfvlHe4TujiYeZlIEEmCJMeVALFlaLpdqkpW09CmCEBWVMZhCdjpJg6aHrTBi7XzTT1uVvW/ey3Vw6h3OFAaoSZCEgggERtz0NELGG3HEZTIIAhM8jpRQO5UfVJSBGm1BbN43L9ytILac4gJAAAjX2z8a35zokyfGiqv2FIV2atyObjh/XVR6RS92BI/Je2yn67m/+dVMEVkZMQn0G5jX5pfwNRMAcFAbUAgDTyke/lVrvj+Y3O/0SvgaiW9uggxpKiNBSDfYBr0a54klsLSo6nYAz7fxrWZxh15AefhttSYCSYDY6k/xrXhzja7VYJAbUpMnWAIPt2H3CueOMFzhB1C1mSEMZwlDcR6x6jaNdYG5iqMmCkP2F0jiJVkWQCkactvDX4VyCVWrslZypPdzDQjaJ8dvbXvBLfEBxnW2kWrb7mXKSUwMs6abajodRyoncWTL9mF3t0sBZGUNJ1VAgkFU7SOWvjVAmyssVu7nDGMJbf4lu2tAuIhLTaiQlSp0EIMdTGmtF8XsbJqxRhnpmVWFpDBccQo6rIU4RH1ZIH+0+dULD2LWywdhmwQGWUtgoE8yNySdTruaknbm/bu8XJSMxS2EetJG865j4bGmxkwgDguqSpQOeRtI0oqfVIAE+NCcJMoeA2zwD7gD91FUgynXSdAOVQVH5P1K/JW0ka5nJiftq60wFeu330B7CQezFrB0zOftFUdI12rFqsd6s+hv5TB4atfZUix1kpuXC0ghDgKwlMHQzyG3PTl471WblJFm/Bk8JXwqU49ndt2FApCddU+sZHUb+UeFMR2wNp30S+caUiGVtGPPNHtmK4JcZS22Hn+GlpBISQVqW4R05xyncmZ1ovg7CWsExKQD9EVa796f3dK83OEILASlsZ4+cKSBKo18etbQLZuSq4QptlTmRspYtgZVrqVLV1Ok9IiIAr9UlKHxeYncpW9pkYYPdbToAB+OVdXMLcygekXGm6Ao9OfXSgd8zb26whKVKVtK1kkiI38qBgt+0eEOtXVnirVy0VlSeOytTqS2QJBR9XTTQHrptSViabRN6sWBdXbD6JbogqERtGgB2rqsJZlRb0jMM2p3rC5GYZU6FJ123OlARwSYdTOpUPZpRlTakOg6ba+NCsKSFJXA9UpjXwonqnLJB5axpvQVTsGoHsvakfad/aKo6Va7Uv/J+AeytqRtxHdf+RVMGUVi9Vlu0j0N+D/Vq+FRlpKinMsTl0BUeX4NWa8AFk/H9kr4VG0LUlod471cSmbDglWC36ZBcys+/NpNEwtJQocMrzHqnrPOKE4Gf5OxHQd1LBGnUz++t7KypBBiJ5VpGda/WUthQLYUqVQBtAHvIpObs3L65zOaJB3KoGnWnV0AJOkyCNfZQjEAEejtoSEhyMxA160QnY08o3yWkZFBDWsDU9KHpbUuM0acxvy8K1qUp19RWZlQTFc7bvukK2SvKPKqohhau4qB3pA8t62KWARmTmjbMdJ8qH4eA3IH1oJk+BralUpVIBg/wqCr/ACfLB7KWu3ru/tFUfKgKXPk+/oraD+86f+xVH1EzXnlfVkf/2Q==" data-csiid="s_2IZ-nYB7ixhbIPoY7DiQc_14" data-atf="4"></div></a></div></div><div class="kb0PBd ieodic" style="grid-area:nke7rc" data-snf="nke7rc" data-sncf="1,2"><div class="VwiC3b yXK7lf p4wth r025kc hJNv6b"><span><em>Hello</em> is a salutation or greeting in the English language. It is first attested in writing from 1826. The greeting "<em>Hello</em>" became associated with telephones&nbsp;...</span></div></div><div class="kb0PBd ieodic" style="grid-area:gdePb" data-snf="gdePb" data-sncf="3"><div class="HiHjCd wHYlTd">‎<a href="https://en.wikipedia.org/wiki/%22Hello,_World!%22_program" ping="/url?sa=t&amp;source=web&amp;rct=j&amp;opi=89978449&amp;url=https://en.wikipedia.org/wiki/%2522Hello,_World!%2522_program&amp;ved=2ahUKEwip_pmJofqKAxW4WEEAHSHHMHEQ0gIoAHoECFwQAQ">"Hello, World!" program</a> · ‎<a href="https://en.wikipedia.org/wiki/Hello_(disambiguation)" ping="/url?sa=t&amp;source=web&amp;rct=j&amp;opi=89978449&amp;url=https://en.wikipedia.org/wiki/Hello_(disambiguation)&amp;ved=2ahUKEwip_pmJofqKAxW4WEEAHSHHMHEQ0gIoAXoECFwQAg">Hello (disambiguation)</a> · ‎<a href="https://en.wikipedia.org/wiki/Hello_Girls" ping="/url?sa=t&amp;source=web&amp;rct=j&amp;opi=89978449&amp;url=https://en.wikipedia.org/wiki/Hello_Girls&amp;ved=2ahUKEwip_pmJofqKAxW4WEEAHSHHMHEQ0gIoAnoECFwQAw">Hello Girls</a> · ‎<a href="https://en.wikipedia.org/wiki/World_Hello_Day" ping="/url?sa=t&amp;source=web&amp;rct=j&amp;opi=89978449&amp;url=https://en.wikipedia.org/wiki/World_Hello_Day&amp;ved=2ahUKEwip_pmJofqKAxW4WEEAHSHHMHEQ0gIoA3oECFwQBA">World Hello Day</a></div></div></div></div></div></div>`

	expectedResult := []GoogleSearchResult{
		{
			URL:         "https://en.wikipedia.org/wiki/Hello",
			Description: "Hello is a salutation or greeting in the English language. It is first attested in writing from 1826. The greeting \"Hello\" became associated with telephones ...",
		},
	}

	results, err := parseGoogleSearchResults(&html)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", results)

	// Clean the actual results
	for i := range results {
		results[i].Description = cleanString(results[i].Description)
	}

	assert.Equal(t, expectedResult, results)
}

func TestCallGoogleSearchResultsTool(t *testing.T) {
	tool := NewGoogleSearchResultsTool(googleSearch.NewRodClient())
	result, err := tool.GoogleSearchResultsFor("what is the weather in japan")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*result)
}

func TestCall(t *testing.T) {
	if os.Getenv("CICD") == "True" {
		t.Skip("Skipping test due to CICD")
	}

	tool := NewGoogleSearchResultsTool(googleSearch.NewNonHeadlessRodClient())
	result, err := tool.GoogleSearchResultsFor("usd to eur prices")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*result)
}
