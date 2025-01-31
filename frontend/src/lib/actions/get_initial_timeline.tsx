"use client";
import useSWR from "swr";
import { GetCollectionOfPostsBySpecificUserAndUsersTheyFollowResponse } from "./get-collection-of-posts-by-specific-user-and-users-they-follow";

type ApiResponse = GetCollectionOfPostsBySpecificUserAndUsersTheyFollowResponse;

type ApiError = {
  message: string;
  status: number;
};

export const getInitialTimeline = () => {
  const fetcher = async (url: string) => {
    console.log('NEXT_PUBLIC_LOCAL_API_BASE_URL:', process.env.NEXT_PUBLIC_LOCAL_API_BASE_URL);
    console.log('NEXT_PUBLIC_USER_ID:', process.env.NEXT_PUBLIC_USER_ID);
    const response = await fetch(url, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const error: ApiError = {
        message: `An error occurred: ${response.statusText}`,
        status: response.status,
      };
      throw error;
    }

    const data = await response.json();
    console.log(`data:${data}`)
    if (!data) {
      return [];
    }
    return data;
  };

  const { data, error } = useSWR<ApiResponse, ApiError>(
    `${process.env.NEXT_PUBLIC_LOCAL_API_BASE_URL}/api/${process.env.NEXT_PUBLIC_USER_ID}/polling?event_type=TimelineAccessed`,
    fetcher,
    {
      revalidateOnFocus: false,
      revalidateOnReconnect: false,
    }
  );

  return { data, error };
};
