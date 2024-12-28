import { PasswordCheckListProps } from "@/types/password.d";
import { useEffect, useState } from "react";

export function PasswordCheckList({
  className,
  value,
  onChange,
  rules,
}: PasswordCheckListProps) {
  const [passwordChecked, setPasswordChecked] = useState(rules);

  useEffect(() => {
    setPasswordChecked((prev) => {
      // @ts-ignore
      return prev.map((rule) => {
        const regex = new RegExp(rule.pattern);
        rule.checked = false;
        if (regex.exec(value)) {
          rule.checked = true;
        }
        return rule;
      });
    });

    // @ts-ignore
    Promise.all(passwordChecked.map((p) => p.checked))
      // @ts-ignore
      .then((result: boolean[]) => {
        let count = 0;
        result.map((b) => {
          if (b) count++;
        });
        return result.length == count;
      })
      .then(onChange);
  }, [value]);

  const checkedIcon = (
    <svg
      className="w-[20px] h-[20px] text-emerald-700  mr-2"
      aria-hidden="true"
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      fill="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        fillRule="evenodd"
        d="M2 12C2 6.477 6.477 2 12 2s10 4.477 10 10-4.477 10-10 10S2 17.523 2 12Zm13.707-1.293a1 1 0 0 0-1.414-1.414L11 12.586l-1.793-1.793a1 1 0 0 0-1.414 1.414l2.5 2.5a1 1 0 0 0 1.414 0l4-4Z"
        clipRule="evenodd"
      />
    </svg>
  );

  const warningIcon = (
    <svg
      className="w-[20px] h-[20px] text-red-700  mr-2"
      aria-hidden="true"
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      fill="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        fillRule="evenodd"
        d="M2 12C2 6.477 6.477 2 12 2s10 4.477 10 10-4.477 10-10 10S2 17.523 2 12Zm7.707-3.707a1 1 0 0 0-1.414 1.414L10.586 12l-2.293 2.293a1 1 0 1 0 1.414 1.414L12 13.414l2.293 2.293a1 1 0 0 0 1.414-1.414L13.414 12l2.293-2.293a1 1 0 0 0-1.414-1.414L12 10.586 9.707 8.293Z"
        clipRule="evenodd"
      />
    </svg>
  );

  return (
    <ul className={className}>
      {/* @ts-ignore */}
      {passwordChecked.map((rule: PasswordRule, i: number) => {
        return (
          <li key={i} className="flex text-sm text-gray-600 my-1">
            {rule.checked ? checkedIcon : warningIcon}
            {rule.label}
          </li>
        );
      })}
    </ul>
  );
}
