package notion

import (
	"context"
	"os"

	openai "audio-to-notion/shared/openai"
	util "audio-to-notion/shared/utils"

	"github.com/dstotijn/go-notion"
)

type NotionClient struct {
	*notion.Client
}

func NewNotionClient() *NotionClient {
	client := notion.NewClient(os.Getenv("NOTION_TOKEN"))
	return &NotionClient{client}
}

func (n *NotionClient) ToBulletedListItemBlocks(arr []string) []notion.BulletedListItemBlock {
	var bulletedListItemBlock []notion.BulletedListItemBlock = []notion.BulletedListItemBlock{}
	for _, item := range arr {
		bulletedListItemBlock = append(bulletedListItemBlock, notion.BulletedListItemBlock{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{Content: item},
				},
			},
		})
	}
	return bulletedListItemBlock
}

func (n *NotionClient) GetBulletedListItemChildren(bulletPoints []string, heading string) []notion.Block {
	var blocks []notion.Block = []notion.Block{}
	blocks = append(blocks, notion.Heading3Block{
		RichText: []notion.RichText{
			{
				Text: &notion.Text{
					Content: heading,
				},
			},
		},
	})
	bulletPointsBlocks := n.ToBulletedListItemBlocks(bulletPoints)
	for _, bulletPoint := range bulletPointsBlocks {
		blocks = append(blocks, bulletPoint)
	}
	return blocks
}

func (n *NotionClient) GetSummaryBlocks(summary string) []notion.Block {
	splitSummary := util.SplitByNumberOfCharacters(summary, 800)
	blocks := []notion.Block{
		notion.Heading3Block{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: "Summary",
					},
				},
			},
		},
	}

	for _, summary := range splitSummary {
		blocks = append(blocks, notion.ParagraphBlock{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: summary,
					},
				},
			},
		})
	}

	return blocks
}

func (n *NotionClient) CreatePage(ctx context.Context, databaseId string, content openai.Completion) (notion.Page, error) {
	databaseProperties := make(notion.DatabasePageProperties)
	databaseProperties["Title"] = notion.DatabasePageProperty{
		Title: []notion.RichText{{
			Text: &notion.Text{
				Content: content.Title,
			},
		}},
	}

	params := notion.CreatePageParams{
		ParentType:             notion.ParentTypeDatabase,
		ParentID:               databaseId,
		DatabasePageProperties: &databaseProperties,
		Children:               []notion.Block{},
	}

	summary := n.GetSummaryBlocks(content.Summary)
	params.Children = append(params.Children, summary...)

	mainPoints := n.GetBulletedListItemChildren(content.MainPoints, "Main Points")
	params.Children = append(params.Children, mainPoints...)

	actionItems := n.GetBulletedListItemChildren(content.ActionItems, "Action Items")
	params.Children = append(params.Children, actionItems...)

	followUps := n.GetBulletedListItemChildren(content.FollowUp, "Follow Ups")
	params.Children = append(params.Children, followUps...)

	arguments := n.GetBulletedListItemChildren(content.Arguments, "Arguments")
	params.Children = append(params.Children, arguments...)

	relatedTopics := n.GetBulletedListItemChildren(content.RelatedTopics, "Related Topics")
	params.Children = append(params.Children, relatedTopics...)

	return n.Client.CreatePage(ctx, params)
}
