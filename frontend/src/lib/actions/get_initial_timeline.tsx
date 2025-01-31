import { Post } from "../models/post";

type InitialTimelineResponse = Post[];

export async function getInitialTimeline(): Promise<InitialTimelineResponse> {
  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_LOCAL_API_BASE_URL}/api/${process.env.NEXT_PUBLIC_USER_ID}/polling?event_type=TimelineAccessed`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },

        cache: "no-store",
      }
    );
    if (response.ok) {
      const responseData: InitialTimelineResponse = await response.json();
      if (!responseData) {
        return [];
      }
      return responseData;
    } else {
      throw new Error(
        `Failed to find post: ${response.status} ${response.statusText}`
      );
    }
  } catch (error) {
    throw new Error("Unable to find post. Please try again later.");
  }
}
